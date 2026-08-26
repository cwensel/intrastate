Model: claude-sonnet-5

# Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0011:C1` ("The kernel change is confined to the match seam... computed from the data TagSet already carries") | Kernel rewrite is larger than claimed: `gate()` only ever walked `row.Guard`; folding match atoms into the same `UndecidedRow.Atoms` payload requires either routing `Match []Tag` through `evaluateAtoms` (which deliberately fails closed on non-guard blocks per `0007:C6`'s fencing) or threading a second undecided-atom channel through `Resolve`, `gate`, and `escapeOrRefuse`. The "confined to the match seam" framing under-scopes the diff and the review effort. | Phase 1 slips; reviewers approve a change that reads as a one-function edit and then discover mid-implementation that `gate`'s signature, `evaluateAtoms`'s block-fencing invariant, and the candidate-loop/gate boundary all move. | §2, premortem |
| C-2 | `0011:C1` ("Entries are deduplicated by `{key, reason}` pair... a key that is both an undecided match atom and an unestablished owned key reports once per distinct reason") | `unknown` folds three semantically distinct sources (kernel match facts, CLI-computed absent-owned-keys, CLI-computed un-run-gate-ids) into one flat list distinguished only by a `reason` string. A consumer that wants to answer "is this row blocked by state I could supply right now, or by a gate I haven't run" has to pattern-match on `reason` values that were never designed as a discriminated union at the type level. | Skill authors write brittle string-matching logic against `reason` (`"absent"`, `"uncomparable"`, `"not-evaluated"`) instead of a typed field, and that logic breaks the next time a fourth reason class is added. | §1, §3 |
| C-3 | `0011:C1` ("a CONFLICTED match key on a row that is also missing owned state... MUST be reported with reason `absent` rather than omitted") | The RDR itself documents a real semantic mislabeling: a conflicted key is reported as `absent` in the compound case, even though the row's remedy (deduplicate readers) is different from an absent key's remedy (bind a reader / supply a tag). The RDR calls this "a reporting-fidelity limit... not a silent drop" but it is still a caller being told the wrong remedy. | A caller acts on the `absent` reason by supplying `--tag key=value`, which does not fix a conflicted-reader problem, and the row stays undecided after the caller does exactly what the payload told them to do. | §1, §3, AT-6 |
| C-4 | `0011:§Consequences` ("an out-of-repo consumer parsing the old field must change even if it adopts `--all`... a caller may owe both [migrations]") | The RDR ships two independent breaking payload changes in one release: the candidate predicate flips (0005→C1 default) AND the field shape/name changes (`unresolved []string` → `unknown []{key,reason}`) in both modes, including `--all`. A caller who adds `--all` to preserve the 0005 candidate SET still breaks on the SHAPE change with no flag to opt out of that half. | Every external consumer of `flow next` breaks on upgrade regardless of which flag they pass, and the "restores 0005's enumeration" framing of `--all` (§Approach, C2) is only half true — it restores the set, not the wire shape, and nothing in the CLI signals that to a caller reading `--all`'s help text. | §2, premortem |
| C-5 | `0011:A6`, `0011:S1`, MVV step 3 (the flagship 21→3 narrowing example) | The one worked example that motivates the whole RDR (`stage=resolved` narrowing 21→3) is a case where the match key is *owned and always served by an invoked reader*. The RDR's own A5 spike shows the opposite case — 0010's decision-table class with no `--tag` — where narrowing does NOT happen and every row stays a candidate with `unknown` entries. The document sells the reader-narrowing win in its headline scenario while conceding in the fine print (A5, §Problem Statement "second case... is the intended behaviour, not a residual defect") that the common unbound/manual case reproduces the original wall-of-candidates symptom verbatim, just with labels attached. | A skill author whose model has manually-supplied (not reader-served) match keys sees the exact same "every row is a candidate" output as before RDR 0011 shipped, and has to read `unknown` entries on 21+ rows to find the one fact they're missing — worse ergonomically than before, since now there's per-row noise instead of a flat list. | premortem, §3 |
| C-6 | `0011:§Approach` ("`flow next` still never selects... a caller facing several candidates is being told the state genuinely admits several") and `0011:Problem Statement` ("Success is therefore 'no row is listed that the supplied state excludes', not a target cardinality") | The RDR defines success as a negative property (nothing wrongly excluded) and explicitly disclaims narrowing to a target cardinality. But the Problem Statement that motivates the whole RDR is a caller who "cannot constrain a skill's next step from" a 21-row list. If the common case (A5: unbound decision-table keys) still returns the full row count, the RDR has not actually solved the stated problem for that caller — it has redefined success so that the original symptom no longer counts as a failure. | A skill author reads the RDR's problem statement, expects `flow next` to now tell them what to do, upgrades, and still gets an unusably long candidate list on any model where match keys aren't all reader-bound — with no signal in the CLI that this is "working as intended" rather than a residual bug. | §3, premortem |
| C-7 | `0011:C1` mini-checks `disposition` table, row "match key present, conflicted" → `{key, uncomparable}`, contrasted with C3's fixture requirement | The RDR requires the conflicted-key fixture to be built from "two invoked readers writing the same owned key with differing values" and explicitly forbids a repeated `--tag` from reaching this path (`parseTags` refuses it upstream). This means the discriminating fixture for the single most novel, most load-bearing behavior in the RDR (present-but-conflicted → reported, not excluded) depends on wiring two readers with overlapping owned-key writes into a test model — a fixture shape that does not exist anywhere in the current fixture set and that the RDR does not show being built, only specifies as a requirement. | Implementation either skips this fixture under time pressure (the "green suite is not evidence" warning in C3/S6 becomes true in the bad direction — a green suite that never actually built the conflicted-reader fixture) or builds a fixture that silently tests the wrong thing (as the RDR itself warns a `--tag`-based attempt would). | §5 AT-6, premortem |
| C-8 | `0011:§Decision Rationale` ("The cost is a larger blast radius than the Stage-2 estimate — recorded in the `large` Profile") and `0011:A7` (Status: Pending) | Two of the eight Critical Assumptions (A7, A8) — precisely the ones that scope the blast radius on `flow resolve`, the verb this RDR admits carries "the largest behaviour change in the RDR" (§Consequences) — are marked `Pending`, to be verified at Stage 6, not before lock. The RDR locks (or is being reviewed for lock) with its own stated riskiest claim unverified, while five lower-stakes assumptions (A1-A6) are fully `Verified`. | Stage 6 finds the oracle census was incomplete (the RDR's own A7 evidence already documents one instance of this happening: "iter-2 already moved that floor from four to five by finding a match-block oracle in a file the first census scoped out"), and a sixth or seventh `internal/resolve` oracle goes red post-lock, requiring a scope amendment after implementation has started. | §1, §2 |
| C-9 | `0011:§Approach` / `0011:C1` (three-valued match verdict lands in the kernel) and `0011:§Decision Rationale` (cost paragraph enumerating 3 non-test call sites of `Resolve`) | The RDR verifies only that `resolve.Resolve` has three non-test call sites and that `TagSet`/`matches` are unexported so the signature "cannot leave `internal/resolve`." It does not audit whether any *external* tooling depends on `flow resolve`'s current behavior of silently dropping an undecidable-match row into an escapable `no_match` (the exact behavior A7/S5 say changes for 3-5 shipped oracles). The consumer-model repo (rdr#tmxk) that *surfaced* this defect is exactly the kind of downstream consumer whose `flow resolve` call sites are never enumerated in this RDR. | A downstream skill or automation that previously relied on an escape row silently rescuing an undecidable-match plan (the exact pattern `TestAdv0007_3_ConflictedKeyIsNotEscapableByNamingIt` was written against, per the RDR's own citation) now gets a hard `guard_unevaluable` refusal in production where it previously got a plan — discovered only by a downstream break, not by anything in this RDR's test plan. | §1, premortem |
| C-10 | `0011:§Load-Bearing Decisions` — Naming: "the flag is `--all`, boolean" | `--all` is a binary axis (match-conditioned vs guard-conditioned-only) bolted onto a system whose only other reporting axis is `--evaluate-gates`. The RDR's own migration analysis (§Consequences) already needs a matrix of default×`--all` and gates-on×gates-off. The Naming rationale explicitly reasons "If a third mode ever appears the migration is `gh`'s own... promote to a `StringEnumFlag`" — i.e., the RDR anticipates its own flag design going stale within one more contract, which is a strong signal that the CLI surface is under-designed for the actual behavior space it's already carving up. | The next RDR that needs a third view mode (e.g., "candidates the guard alone excludes, independent of match, with gates evaluated") has to choose between a breaking flag migration or a second boolean flag, and `flow next --help` accretes flags rather than composing them. | §2 |

# 1. The three most likely ways implementation goes wrong

## Failure 1: the kernel change is bigger than "confined to the match seam" implies, and Phase 1 blows its estimate

**Root cause.** The Approach section states plainly: "The kernel change is confined to the match seam: `TagSet.matches`' `bool` becomes a three-valued verdict plus the atoms that could not be decided, computed from the data `TagSet` already carries" (`0011:§Technical Design`, and repeated almost verbatim at `0011:C1`'s prose lead-in). This framing, and the Implementation Plan's Phase 1 ("Intent: give `internal/resolve`'s match comparison a three-valued verdict... reusing `0007:C8`'s `Reason` vocabulary"), reads as a single-function signature change plus payload reuse.

**The specific passage that enabled it.** Grounding against the actual source (`internal/resolve/resolve.go`) shows the real shape of the change:

- `Resolve`'s candidate loop currently does `if !view.matches(row.Match) { continue }` — a hard prune, before `gate()` is even called.
- `gate(rows []Row, seam GuardEvaluator, view TagSet)` only ever consumes `row.Guard []GuardAtom`. It has no parameter for match atoms and no path for them.
- `evaluateAtoms` (called inside `gate`) explicitly **fails closed** on any block that isn't `BlockAll`/`BlockUnless`: `if !isGuardBlock(atom.Block) { continue }`, with a comment explaining this is deliberate — an atom in `BlockMatch` "is never handed to the seam" and "contributes NO operand."
- `UndecidedRow.Atoms []UndecidedAtom` is populated exclusively inside `evaluateAtoms`/`gate`.

C1 requires an `indeterminate`-match row to (a) *not* be pruned by the candidate-loop filter (survive past today's `if !view.matches(...)`), (b) reach `gate()` and be treated as a non-selectable-but-owned-obligation-bearing row, and (c) have its undecided match atoms land in the *same* `Refusal.Undecided` / `UndecidedRow.Atoms` payload that today is guard-only and is deliberately fenced against non-guard blocks by `evaluateAtoms`. That means either widening `isGuardBlock`'s fencing (contradicting the very invariant `0007:C6`/`0007:C8` established, and which the RDR cites approvingly) or building a second, parallel undecided-collection path that has to be merged with the guard one before `gate`'s two return points. Neither is "computed from the data TagSet already carries" in the way the prose implies — both require rewiring the boundary between the candidate loop and `gate`, which today is a hard filter-then-forget, into a filter-that-remembers-what-it-almost-filtered.

The RDR partially acknowledges this with real engineering language ("`0007:C11`'s condition — pruning is safe only on decided atoms", "An `indeterminate` row therefore SURVIVES into `gate` as a non-selectable row") but never states the actual diff shape: does `gate` grow a new parameter? Does `Row` grow a computed field? Is the match verdict computed once in `Resolve` and threaded down? This is exactly the kind of decision an implementer has to invent at code time because the RDR's Technical Design describes data flow at the semantic level ("→ resolve.Resolve over the one-row table → the kernel evaluates each match atom three-valued... ⇒ candidate") without naming the call graph shape that survives the D8/C11 pruning-order constraint.

**Symptom the user sees.** Phase 1 (nominally "kernel, small, isolated") takes materially longer than planned, because the implementer discovers mid-flight that `gate`'s signature must change, `evaluateAtoms`'s deliberate block-fencing has to be either violated or duplicated, and `Resolve`'s candidate loop and `gate`'s owned-state-precedence logic (`owned_state_unavailable` before `guard_unevaluable`, per S7) both need to be re-derived for the three-valued match case rather than just "extended." The Profile is already `large`; this RDR under-estimates even that.

## Failure 2: `unknown`'s three-source fold produces a payload that looks typed but isn't, and consumers build brittle logic against it

**Root cause.** C1 mandates: "Every `indeterminate` match atom, every guard fact the kernel could not decide, every owned key no invoked reader established, and (absent `--evaluate-gates`) every gate id MUST appear in that candidate's `unknown` list as a `{key, reason}` pair." Four distinct *causes* (kernel match-indeterminate, kernel guard-indeterminate, CLI-computed missing-owned, CLI-computed un-run-gate) are collapsed into one list keyed by `{key, reason}`, where `reason` is one of three string tokens (`absent`, `uncomparable`, `not-evaluated`) that do not map 1:1 to the four causes — `absent` alone covers both "match key absent from view" and "owned key no reader established" (per the `disposition` table: both report `{key, absent}`), and, per C1's own compound-case carve-out, sometimes covers "match key conflicted" too (see Failure/Finding C-3 below).

**The specific passage that enabled it.** The `disposition` table (`0011:C1` mini-checks) is the tell: five different input classes ("match key absent from view", "owned key no reader established", the compound-case conflicted-key special case) all resolve to the identical wire tuple `{key, absent}`. A consumer cannot distinguish "bind a reader for this owned key" from "supply this match tag" from "your two readers conflict and I'm calling it absent because I can't say otherwise" without re-deriving side-channel knowledge the payload does not carry (e.g., cross-referencing `Required`/owned key names against the model, which the RDR does document as available data, but never states as the intended consumption pattern).

**Symptom the user sees.** A skill author or downstream tool writes `if reason == "absent": bind_the_reader()` logic, and that logic silently mis-fires on the compound conflicted-key case (does nothing useful, since there's no reader to bind — the fix is deduplicating owned-key-writing readers) and provides no signal that a different remedy applies. The RDR calls this "a reporting-fidelity limit of the compound case, not a silent drop" (`0011:C1`) — true in the sense that a key name appears — but false in the sense that matters to an automated consumer: the *reason* is wrong, not absent.

## Failure 3: the flagship narrowing story (21→3) does not generalize to the RDR's own second worked example (0010's decision tables), and nobody notices until a real skill hits it

**Root cause.** The Problem Statement, Approach, and MVV all lead with the `models/rdr.toml` `stage=resolved` example: 21 candidates narrow to 3 because `stage` is `required = true` and always served by the `rdr-status` reader (A6). This is the RDR's proof that the design works. But A5 — verified via spike, not hand-waved — shows that over 0010's decision-table class with **no `--tag` supplied**, "every ordinary rule" stays a candidate, and each rule's match keys land in `unknown`. The Problem Statement itself concedes this is "the old symptom's shape with a diagnosis attached" and asserts, without much argument, that it "is the intended behaviour, not a residual defect."

**The specific passage that enabled it.** `0011:§Problem Statement`: "where it is observed and unsupplied (0010's decision-table class with no `--tag`, A5) every row remains a candidate and each names the undecided key under `unknown`. That second case is the old symptom's shape with a diagnosis attached, and it is the intended behaviour, not a residual defect." This is a definitional move, not an engineering one: the RDR redefines "success" as "narrows only when the state is knowable" (§Problem Statement: "Success is therefore 'no row is listed that the supplied state excludes', not a target cardinality") specifically to make the A5 case not count as failure. But the caller who filed kata `1mv1` didn't file it because "no row is listed that the state excludes" was violated abstractly — they filed it because they got a 21-row wall and couldn't use it. A caller on an unbound decision-table model gets that same 21+-row wall (now with row-level annotations) and the RDR has pre-emptively defined that outcome as correct.

**Symptom the user sees.** A skill built against a decision-table-shaped model (exactly 0010's class, exactly the class the joint-check ties this RDR to) runs `flow next` with partial or no `--tag` input and gets the full row count back, now with `unknown` entries attached to every row — more bytes on the wire, same practical unusability, and the caller has no CLI signal distinguishing "this model just can't narrow" from "you forgot a flag."

# 2. The one section that will be rewritten within 6 weeks of shipping

**`0011:C1`'s "compound row" clause** (the paragraph beginning "One disposition carries only part of the row's facts... For that row the `unknown` list MUST carry the missing owned keys from `MissingOwned` plus the row's match-block atoms over keys ABSENT from the assembled view... a CONFLICTED match key on a row that is also missing owned state... MUST be reported with reason `absent` rather than omitted").

This is the section that will be rewritten, for three converging reasons:

1. **It is the one place the RDR admits its own vocabulary is inadequate.** The RDR built exactly two reason tokens for the match/guard seam (`absent`, `uncomparable`, reusing `0007:C8`) plus one CLI-only token (`not-evaluated`), and then hits a case — conflicted-and-also-missing-owned — where the honest answer isn't representable, and settles for mislabeling a conflicted key as `absent`. The RDR itself flags this as fragile: "This is a reporting-fidelity limit of the compound case, not a silent drop, and it is the ONLY case in which `unknown` does not distinguish `uncomparable` from `absent`." A design that needs a footnote explaining why its own invariant ("uncomparable vs absent, always distinguished") has exactly one carve-out is a design under strain at that seam.

2. **It is reachable and will be reported.** This isn't a theoretical edge case — the RDR's own S7 test scenario exists specifically to exercise it, meaning it's expected to happen in practice (a row with both an owned-state gap and a conflicted match key, e.g., two readers disagreeing on a key that's also `required` by a different clause of the row). Once a real skill author hits `{key, absent}` on a genuinely conflicted key, follows the "bind a reader / supply the tag" remedy the `absent` reason implies (per `0011:§Consequences`: "`absent` and `uncomparable` have different remedies... and the payload names which"), and it doesn't fix anything, that's a filed defect against the payload's own documented promise.

3. **The root cause — `gate` returning `owned_state_unavailable` with an empty `Undecided` before the guard/match-undecided payload is even collected — is a precedence choice the RDR preserves rather than re-derives** ("the existing precedence is preserved, not reordered"). That precedence was set by RDR 0007 for a two-source payload (owned vs. guard-undecided). This RDR adds a third source (match-undecided) into the same precedence stack without revisiting whether "owned state wins, full stop" is still the right precedence once match facts are in play — it just inherits the ordering and patches around its consequence with a mislabeled reason token. The natural fix once this is live — carry `Undecided` alongside `MissingOwned` on the same refusal instead of discarding it — is a `Refusal` shape change, which is exactly the kind of follow-up RDR this critique predicts within 6 weeks.

# 3. The one assumption that will not survive first contact with a real user

**`0011:§Problem Statement`'s reframe of "success":** *"Success is therefore 'no row is listed that the supplied state excludes', not a target cardinality."*

This is the assumption load-bearing enough to be stated as doctrine, and it is the one a real user will reject on contact. The actual human (or skill) driving `flow next` does not care, and was never asking, whether the list is *logically sound* relative to what's knowable — they were asking a *practical* question: "what should I do next." The kata that started this whole RDR (`1mv1`) was filed because the list was too long to be useful, not because it contained a logically indefensible row. RDR 0011 fixes the logical-soundness half rigorously (three-valued match, kernel-verified, extensively cross-checked against XACML/PTaCL/SQL NULL precedent) and, in the one place its own evidence (A5) shows the practical-usability half is not fixed — decision-table-shaped or otherwise under-bound models — declares that case out of scope by redefinition rather than by argument tied to the caller's actual need.

A real user with a partially-bound model (which, per the RDR's own joint-check with 0010, is not a corner case but an entire *class* of model this system supports) will run `flow next`, get the same wall of rows they got before, and report "this didn't fix it" — technically wrong under the RDR's redefined success criterion, and correct under the criterion the kata was actually filed against.

# 4. The premortem

*Six weeks post-ship. Two independent incidents, both traced back to RDR 0011.*

**Incident A — "the narrowing doesn't narrow."** A skill author building against a newly-modeled decision-table workflow (structurally identical to 0010's class — multiple independent match dimensions, several owned tags, no single dominant reader) runs `flow next --model workflow.toml --artifact state=<partial>`. They get back all 14 ordinary rows as candidates, each carrying two or three `{key, absent}` entries in `unknown`. They file a bug: "RDR 0011 was supposed to fix this." Support (or the on-call maintainer) closes it "working as designed," citing `0011:§Problem Statement`'s success redefinition and A5's spike evidence. The bug reporter is unsatisfied — from their seat, `flow next` still returns the whole alphabet on a decision-table model, exactly the symptom the RDR's own headline claimed to have fixed, and the fix that DOES apply (bind the readers, or supply every `--tag`) is exactly as much manual work as reading the whole table by hand was before. `internal/cli/flow_next.go::summarize` is doing more work — walking `Row.Atoms` by `Block`, reading `Refusal.Undecided` off the kernel probe, deduping by `{key, reason}` — to produce a list that is, for this caller's model shape, functionally the same size it was pre-RDR-0011.

**Incident B — "the conflicted-key remedy doesn't work."** A different skill, driven by two readers that both write an `owned` `env` tag (one from a Kubernetes context reader, one from a local `.env` file reader, disagreeing during a migration window), runs `flow next`. Rows that match on `env` come back as candidates with `{"env", "absent"}` in `unknown` — because this row also happens to be missing a different owned key, tripping C1's compound-row carve-out, which the RDR itself specifies should report the conflicted `env` key with reason `absent` rather than `uncomparable`. The caller, trusting the payload (per C3's help text: "a match or guard key that is unsupplied... leaves a row a candidate with the key listed under `unknown` and its reason named"), does what `absent` tells them to do: they bind another reader for `env`, or pass `--tag env=prod`. Neither fixes anything, because the actual problem — the *existing* two readers disagree — is invisible in this payload shape. The caller escalates as "the fix doesn't fix it," and whoever triages finds `0011:C1`'s own text pre-committing to exactly this outcome: "the missing owned key is the actionable remedy the caller acts on first" — an editorial decision that the *other* (wrong) remedy is what the caller sees, made at RDR-review time and never revisited against a live conflicted-reader incident.

Both incidents trace to the same root: the RDR is verified rigorously against **what the kernel can prove**, and comparatively thin against **what the caller, in the moment of receiving a payload, will conclude and do**. `flow_next.go::runFlowNext` and `summarize` will ship correct relative to every one of C1/C2/C3's normative clauses, and relative to A1-A6's verified evidence — and still leave two classes of real caller no better off, or actively misdirected, than before the RDR shipped.

# 5. Acceptance tests that would have caught each failure at RDR-review time

**AT-1 (catches Failure 1 — kernel scope underestimate).**
```
Given the Implementation Plan states Phase 1 is "confined to the match seam"
When the reviewer traces the call graph required for C1's compound-row and
  selection-site clauses through Resolve -> gate -> evaluateAtoms
Then the RDR MUST state explicitly whether gate()'s signature changes,
  whether isGuardBlock's block-fencing is widened or bypassed, and whether
  a new field is added to Row or Refusal to carry match-vs-guard undecided
  atoms separately before being merged
And if none of these is named, the RDR MUST NOT claim the change is
  "confined to the match seam" in the Approach section
```

**AT-2 (catches Failure 2 — `unknown` payload ambiguity).**
```
Given C1's disposition table shows "match key absent from view" and "owned
  key no reader established" both reporting {key, absent}
When a consumer needs to choose a remedy from a single unknown[] entry
Then the RDR MUST demonstrate (via a fixture or a written trace) that a
  consumer CAN distinguish "supply this match tag" from "bind a reader for
  this owned key" from the wire payload alone, without side-channel model
  knowledge
And if it cannot, C1 MUST add a fourth field (a cause discriminator, or
  split unknown into typed sub-lists) rather than relying on the reader to
  infer cause from key name
```

**AT-3 (catches Failure 3 / Assumption in §3 — narrowing doesn't generalize).**
```
Given the RDR's flagship example (A6, models/rdr.toml) narrows 21->3 because
  every match key is reader-owned
And A5's spike shows 0010's decision-table class narrows 0 rows on no --tag
When the Problem Statement is reviewed
Then the RDR MUST report, alongside the 21->3 number, the SAME metric
  (rows-before / rows-after) for at least one decision-table-shaped model
  representative of the 0010 class this RDR cross-cites
And if that number is N/N (no narrowing), the RDR MUST NOT present the
  21->3 case as evidence the general problem is solved without an equally
  prominent statement of the class for which it is not
```

**AT-4 (catches the "success" redefinition in §3).**
```
Given the Problem Statement declares success as "no row is listed that the
  supplied state excludes, not a target cardinality"
When a reader compares this to the original kata (1mv1)'s complaint
Then the RDR MUST show that the redefined success criterion, applied to
  the ORIGINAL reported model+state that triggered the kata, actually
  produces a materially shorter, more actionable list than before
And if a broad class of models (decision tables, A5) is exempted from
  narrowing by this same redefinition, the RDR MUST name that exemption
  in the Problem Statement itself, not only in a later Critical Assumption
```

**AT-5 (catches the compound-row mislabeling, Incident B).**
```
Given C1's compound-row clause reports a CONFLICTED match key as reason
  "absent" when the row is also missing an owned key
When a caller acts on the "absent" remedy (bind a reader / supply a tag)
  for a key that is actually conflicted (two readers disagree)
Then a test MUST assert that the caller's action (adding --tag env=X)
  does NOT silently succeed in making the row a clean candidate while the
  underlying reader conflict remains -- i.e. the test must show the
  mislabeled remedy visibly fails to resolve the row, and the RDR's
  Failure Modes section MUST document this specific dead-end explicitly
  (not just "reporting-fidelity limit")
```

**AT-6 (catches the conflicted-key fixture risk, C-7).**
```
Given C3 requires the conflicted match-key fixture to be built from two
  invoked readers writing the same owned key with differing values, and
  explicitly forbids a repeated --tag from reaching this path
When the fixture is authored during implementation
Then the PR/commit implementing S3's fourth row MUST include a NEW test
  reader pair (not reuse of an existing single-writer reader) whose
  OwnedSnapshot() outputs are shown to collide on one key
And the review MUST confirm this fixture exists and passes BEFORE citing
  "green 0005/0011 suite" as evidence C1 shipped, per C3's own warning
  that a green suite is not such evidence
```

**AT-7 (catches the two-migration cost in C-4).**
```
Given the payload shape change (unresolved[]string -> unknown[]{key,reason})
  applies in BOTH default and --all modes
When help text or a migration note is reviewed
Then flow next's --all help text and any migration guidance MUST state
  explicitly that --all restores the 0005 candidate SET but NOT the 0005
  payload SHAPE, so a caller relying only on --all still has a breaking
  change to absorb
And if no such statement exists in Short/Long help or docs/cli-output-
  contract.md, C3 MUST be revised to add it
```

**AT-8 (catches the Pending-assumption-at-lock risk, C-8).**
```
Given A7 and A8 are Status: Pending, to be verified at Stage 6, and A7's
  own evidence documents the oracle census already grew from four to five
  between review iterations
When the Finalization Gate's Assumption Verification is completed
Then the gate response MUST NOT treat A7/A8 as low-risk simply because
  they are "to be verified at Stage 6" -- it MUST state what happens to
  C1's lock status if Stage 6 finds a SIXTH internal/resolve oracle (a
  demonstrated failure mode of the census method itself, per A7's own
  "iter-2" note), and whether that triggers a scope amendment or blocks
  merge
```
