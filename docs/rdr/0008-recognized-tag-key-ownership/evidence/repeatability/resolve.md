model: claude-opus-5[1m]

# Repeatability resolve — RDR 0008

Dispositions for `diff.md`. Origin ledger = the diff's D-1…D-6. Grounded
against `internal/resolve/resolve.go` on `main`, `{RDR_RESOURCES}`, and the
RDR's own decided text before any edit.

| # | Disposition | Origin | Section touched |
| --- | --- | --- | --- |
| D-1 | **pin** — advisory payload given two named fields, rule id `reserved-tag-key/near-miss`, and its own non-category channel; trigger disambiguated to **disjunctive** (fold *or* trim) | D-1 | Load-Bearing Decisions / Identity; scenario 4 |
| D-2 | **pin** — predicate reports first breach, single error, scan order latitude; explicitly severed from the "detected whenever present" clause it was being confused with | D-2 | Normative Contracts block 4; scenario 6 |
| D-3 | **pin** — predicate read-domain stated at block 4 where the predicate is defined, not recovered from block 5's A11 aside | D-3 | Normative Contracts block 4 |
| D-4 | **pin** — payload split by direction: `required name` → `remedy name`, second rule id `reserved-tag-key/author-must-rename`, renderer forbidden from presenting the reserved key as required in the rename-away direction | D-4 | Normative Contracts block 3; scenarios 2, 7; Phase 3; Risks; Failure Modes; MVV |
| D-5 | **pin** + new assumption — `Input.Owned`/`Observed` are `[]Tag` slices; duplicate reserved keys admissible; booked as **A13 (Verified, Source Search)** | D-5 (grounding, not run disagreement) | Critical Assumptions (A13); block 4; scenario 6 third variant |
| D-6 | **leave non-normative** | D-6 | none — RDR already states exporting the spelling constant is implementation latitude (Technical Design), and the behavioral-conformance requirement (premortem P-5) makes either choice safe |

## Grounding-gate notes

- **D-3 / D-5 grounded against code**: `Input` carries `Table Table`, so the
  runs' `in.Table.Rows` inference was correct — D-3 is an RDR silence, not a
  wrong claim. `Owned []Tag` / `Observed []Tag` with
  `Tag{Key, Value string}` refutes all three runs' map model, which is what
  promoted D-5 from "irrelevant shape trivia" (the runs' own assessment) to a
  contract edit.
- **No finding was dismissed for citing absent code.** Every cited symbol
  (`Input`, `Tag`, `assemble`, `recognizedTagKey`, `missingOwned`,
  `RequiresOwned`) resolves on `main`.
- **Not re-litigated**: the 0008↔0009 precondition order (A11 / Stage 7.1) and
  the exported predicate's *name* (explicit latitude). Both appeared as GUESS
  clusters; both are already-decided-open in the RDR, so per the grounding gate
  they are re-raises, not defects. run-1 flagged the first defensively as "not a
  defect; recorded so the diff does not count it as one" — correct, and honored.

## Charted

None. Every finding landed inside the identity contract this RDR already owns —
no net-new scope was absorbed, and no successor RDR is owed.

## Needs (re)verification — carried to Stage 6

- **A13 (new, Verified)** — booked as Verified at authoring since the source
  read is direct and current, not deferred. Listed here because Stage 6 should
  confirm the `[]Tag` shape has not moved under RDR 0001 before lock.
- **No previously-Verified assumption was invalidated.** The edits pin
  under-specified surfaces; none contradicts A1–A12's established evidence.
  A7, A9, A10, A11, A12 remain `Pending` on their existing plans, untouched by
  this pass.
- **New load-bearing claims introduced by these fixes**, all normative text
  rather than external-behavior claims, so none needs a spike:
  `reserved-tag-key/author-must-rename` and `reserved-tag-key/near-miss` are
  tokens this RDR authors (like `reserved-tag-key/kernel-owned` before them);
  the disjunctive advisory trigger is a rule this RDR authors; the first-breach
  arity is a rule this RDR authors. Each is pinned by a named scenario
  (4, 6, 7), which is the falsifiability this RDR requires of its own MUSTs.

## Tiebreakers escalated

None. Every fork collapsed on evidence:

- D-1's fold-vs-trim conjunction was decided by **scenario 4's own
  requirements** — it demands the advisory on `Recognized` (fold-only) and
  `" recognized"` (trim-only), so a conjunctive reading is refuted by the
  RDR's existing test matrix rather than by preference.
- D-4's direction split was decided by the **Trade-offs section's own
  mitigation claim** ("Mitigated by the failure *data* — offending name,
  required name, rule identifier"), which misfires for the rename-away rule.
  The RDR was already relying on a payload that gave wrong advice half the
  time; fixing it is repair, not a design call.
- D-2's arity was decided by the **cost asymmetry**: a producer holding a
  programmer mistake is not owed an exhaustive list, and aggregation would
  require specifying scan order the RDR deliberately leaves free.
