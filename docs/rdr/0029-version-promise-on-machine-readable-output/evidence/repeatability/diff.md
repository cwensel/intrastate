Model: claude-sonnet-5

# Repeatability-lite DIFF — RDR 0029, run-1

Comparing `evidence/repeatability/run-1.md` against the RDR's `kind=="C"`
spans (C1-C4), MVV, S1-S9, and Load-Bearing Decisions, read verbatim via
the projector. Widened per run-1's own log into `§illustrative-code` and
`§load-bearing-decisions` to confirm those widenings were warranted, and
into A5/A7/A6 (Critical Assumptions) to adjudicate the GUESS markers,
since C4 explicitly ties two of the six owed seams to A5's Pending status.

---

## Findings

None admitted.

No GUESS in run-1, and no independent read of C1-C4/MVV/S1-S9, surfaced a
step whose order the RDR leaves unfixed, a field whose owner it leaves
unnamed, a transform whose tie-break it leaves unstated, or an error mode
it leaves unspecified. Every candidate traced below resolves to either a
naming question (inadmissible) or an assumption the RDR already tracks
under an explicit Pending id with a stated fallback (not a silence — a
documented open item with an owner and resolution mechanism, per A5/A7's
own "If wrong" clauses).

---

## GUESSes checked and dropped as run artifacts

1. **Exact Go identifier for the `SchemaVersion` constant** — C1 fixes the
   value (`"0.1"`), the type (string, `MAJOR.MINOR`), the home package
   (`clierr`, the leaf), and the cardinality (one constant, read not
   declared by `respond`). It does not spell the identifier. This is
   naming/wording, explicitly excluded by the admissibility rule — dropped,
   not a finding.
2. **`clierr.Codes()` / `graphlint.UnknownReasons()` signature or
   location** — contingent on whether these seams are mechanical at all.
   C4's census marks both "Owed" but flags them non-mechanical, and A5
   (Status: Pending, Method: Spike) is on record as the exact open
   question, with an explicit fallback already written into C4: "the tier
   on that row STANDS and the row records `seam: none (prose-only)`."
   The RDR does not fail to determine the outcome — it names both possible
   outcomes and assigns the choice to A5's spike before lock. Run artifact:
   the run should have widened into A5 and reported this as tracked rather
   than open; the run's own text acknowledges A5 Pending but still marks
   these as GUESSes needing resolution the RDR already structurally
   supplies.
3. **"Plain exported const vs. a zero-arg func" framing (helper #1 in
   run-1 §2)** — collapses into #1: C1 says "constant," not accessor, for
   `schema_version` specifically (the C4 "enumeration seam" pattern of
   exported accessor functions applies to tiered *vocabularies*, a
   different mechanism from the single scalar version marker). No genuine
   shape ambiguity — dropped.
4. **CI seam-diff snapshot file location/format** — A7 (Status: Pending,
   Method: Spike) is the RDR's own record of this exact question,
   enumerating the candidate mechanisms ("a hidden verb, a `go:generate`
   golden file, or a test fixture") and stating the fallback if none
   lands ("the clause reverts to review-only disclosure, and that limit
   is recorded in Capability Dependencies rather than papered over"). Run
   artifact, not a silence — the RDR defers the choice to a named spike
   with a stated consequence either way.
5. **Whether the two A5-contingent seams land as real accessors or
   `seam: none (prose-only)`** — same item as #2; restated by the run as
   a separate GUESS in the pseudo-code block. One drop covers both
   instances.

## Where run and RDR agree

`respond.Types()` returns exactly `{"ok"}`, `"failed"` MUST NOT enter the
seam, and the refusal record carries no `type` key at all — the run's
reconstruction matches C4's census row and S6 verbatim, including the
reasoning (an assertable frozen seam can't pin a member no consumer can
ever observe). This is a contract the run reproduced correctly rather than
guessed at, confirming C4/S6 are determinate here.

---

## Escalation judgement

No load-bearing or cross-RDR contract silence surfaced. The two genuinely
open items (A5's seam mechanics, A7's CI snapshot mechanism) are already
tracked as Pending assumptions with spikes and stated fallbacks inside
0029 itself — they are not gaps in the peer contracts 0029 cites (0002,
0003, 0005, 0006, 0007, JDR 0001); this diff found no daylight between
0029's C1-C4 and what those peers are quoted as requiring (e.g., C4's
verbatim quotation of `0003:C7` on the operator vocabulary, `0005:C1`'s
bare-refusal shape, `0006:C17`'s disposition guarantee, `JDR 0001 §D12`'s
append-only precedent). Escalation to the full x3 lens is NOT warranted on
this pass's evidence.

The run is not shared-model overfit: 7 GUESS markers, explicit widen-log,
and a section (§2 helper #1, the A5 seam GUESSes) where the run flags its
own uncertainty rather than silently picking one outcome. That is the
expected shape for a spike-gated RDR, not a red flag.
