# Cluster Reconcile Report 0002-0009

Date: 2026-08-11

Model: claude-opus-5

## Membership

Cluster = the Final-and-unimplemented RDRs under `docs/rdr`. RDR 0001 is
`Implemented` and is therefore **out of scope** (code is the source of truth);
it appears below only as the seam the newer members bind.

| RDR | Status | Implemented? | Profile | Relatedness |
| --- | --- | --- | --- | --- |
| 0002 | Final | No | large | Transition-model data format; consumed by every other member. |
| 0003 | Final | No | large | Guard operator vocabulary + exhaustiveness; named obligation destination for 0007. |
| 0004 | Final | No | large | Accessor execution safety; obligation destination for 0007 A6b and 0009. |
| 0005 | Final | No | mid | User-facing CLI contract; the surface every refusal must reach. |
| 0006 | Final | No | large | Graph-lint authority; the static half of the static/runtime seam. |
| 0007 | Final | No | foundational | NEW cross-RDR producer at the 0001↔0003 seam (guard evaluation domain). |
| 0008 | Final | No | foundational | NEW cross-RDR producer at the 0001↔0002 seam (recognized tag key name). |
| 0009 | Final | No | foundational | NEW cross-RDR producer at the 0001↔0002 seam (escape-row shape). |

The cluster is real and the gate is correctly triggered: 0007/0008/0009 all
declare `Seam Lineage: area:internal-resolve`, all three landed on 2026-08-09
against peers locked on 2026-06-19, and all three file binding obligations on
peers. This is precisely the "later lock drifts an earlier peer" case no
per-RDR gate covers.

### Pairs run (13)

Trimmed to explicitly interacting pairs named by Predecessors, Peer-RDR
assumption evidence, or shared cross-cutting seam. Pairs among 0002-0006 were
reconciled by the prior `0001-0006` pass and are not re-run except where a new
member reopens them.

`0007-0002`, `0007-0003`, `0007-0004`, `0007-0005`, `0007-0006`, `0007-0008`,
`0007-0009`, `0008-0002`, `0008-0009`, `0008-0005-0006` (combined),
`0009-0002`, `0009-0004`, `0009-0005`.

## The set-level finding that governs this report

Mechanically verified during this pass:

```
grep -cE 'RDR 0007|RDR 0008|RDR 0009' 0002…md 0003…md 0004…md 0005…md 0006…md
→ 0, 0, 0, 0, 0
```

**Not one of the five older members references any of the three new members.**
Citation traffic is entirely one-directional. Meanwhile 0007 names
`Obligation destination (NAMED): RDR 0003's implement stage` three times (A10,
A12, A15), routes A6b to RDR 0004's implement stage, and 0008/0009 route
further obligations to 0002/0004/0005.

RDR 0003's Prerequisites read `- [x] All Critical Assumptions verified`, and
RDR 0002's likewise — both checked, both predating the obligations now bound to
them. An implementer starting from either RDR reads a green checklist that is
materially false.

This is the cluster's defining defect: **the set has been reconciled document by
document, but the obligations flow only one way, into peers whose locked text
has no knowledge of them.** Every pairwise scan below is a symptom of it.

## Pairwise findings table

Severity/ownership as returned by each scan. 31 findings; 24 `blocks-impl`.

| Pair | Finding | TYPE | SEVERITY | Ownership | Disposition |
| --- | --- | --- | --- | --- | --- |
| 0007-0003 | `Row.Guard` is one string/row but 0003's Identity key is per-atom; N-atom rows cannot round-trip | round-trip | blocks-impl | joint | JOINT-DECISION → JD-1 |
| 0007-0003 | `all … exists = false` vs 0003's "Positive guard atoms MUST live in `all`" | contradiction | risks-impl | joint | JOINT-DECISION → JD-1 |
| 0007-0003 | `contains` over absent set-valued tag has no image in 0003's coverage algebra | gap | risks-impl | joint | JOINT-DECISION → JD-4 |
| 0007-0002 | 0002's match-count resolver flow vs 0007's pre-count aggregation veto | contradiction | blocks-impl | joint | JOINT-DECISION → JD-2 |
| 0007-0002 | Guarded escape rows: shape 0007 quantifies over is ungranted by 0002 | gap | blocks-impl | joint | JOINT-DECISION → JD-2 |
| 0007-0002 | `RequiresOwned` absent from 0002's round-trip preservation set | round-trip | blocks-impl | joint | JOINT-DECISION → JD-3 |
| 0007-0008 | Absent `recognized` key: 0007 says unevaluable (table-wide veto), 0008/0002 imply plain no-match | contradiction | blocks-impl | joint | JOINT-DECISION → JD-4 |
| 0007-0008 | `exists` on the reserved key is an unowned outcome-gating backdoor | gap | blocks-impl | joint | JOINT-DECISION → JD-4 |
| 0007-0009 | Two whole-table entry preconditions on one `Resolve`, no precedence rule | gap | blocks-impl | joint | JOINT-DECISION → JD-5 |
| 0007-0009 | `RequiresOwned` on escape rows: 0007 defines by `Writes`, 0009 empties `Writes` | round-trip | blocks-impl | joint | JOINT-DECISION → JD-3 |
| 0007-0009 | Both claim RDR 0001's programmer-mistake channel — error vs panic | contradiction | risks-impl | joint | JOINT-DECISION → JD-6 |
| 0007-0004 | A6b read-completeness obligation binds nobody (0004 mentions 0007 zero times) | gap | blocks-impl | joint | JOINT-DECISION → JD-7 |
| 0007-0004 | Truncated accessor reads round-trip green | round-trip | blocks-impl | joint | JOINT-DECISION → JD-7 |
| 0007-0005 | `owned_state_unavailable` has no CLI code row in 0005's twelve-row table | gap | blocks-impl | joint | JOINT-DECISION → JD-8 |
| 0007-0005 | `flow-guard-unevaluable` scoped to "supplied facts" vs 0007's provenance-blind view | contradiction | blocks-impl | joint | JOINT-DECISION → JD-8 |
| 0007-0005 | Envelope cannot carry `Refusal.Rows`/`MissingOwned` ("no new envelope fields") | gap | blocks-impl | joint | JOINT-DECISION → JD-8 |
| 0007-0005 | `read-state`→`--tag` loses provenance, defeating non-escapability | round-trip | blocks-impl | joint | JOINT-DECISION → JD-9 |
| 0007-0006 | Lint can pass a table runtime refuses `guard_unevaluable`; 0006 silent | contradiction | blocks-impl | joint | JOINT-DECISION → JD-4 |
| 0007-0006 | No lint category for "guard may be unevaluable at runtime" | gap | risks-impl | joint | JOINT-DECISION → JD-4 |
| 0008-0002 | Empty-outcome input: declaration total, emit conditional — key absent vs declared-empty | round-trip | blocks-impl | joint | JOINT-DECISION → JD-10 |
| 0008-0002 | "Narrows nothing" false: name constraint invalidates 0002's three canonical fixtures | contradiction | blocks-impl | joint | JOINT-DECISION → JD-10 |
| 0008-0002 | No pointer from 0002 to 0008; handoff form undetermined | gap | blocks-impl | joint | JOINT-DECISION → JD-10 |
| 0008-0009 | **0008 A6/A11 `Verified` on a false negative-existential; citations stale** | contradiction | blocks-impl | **single-RDR (0008)** | **SPEC-DEFECT → 0008 to Draft** |
| 0008-0009 | Two entry preconditions, no precedence; oracle base commit unnamed | gap | blocks-impl | joint | JOINT-DECISION → JD-5 |
| 0008-0009 | 0009's frozen-suite oracle not invariant under 0008's fixture rename | round-trip | risks-impl | joint | JOINT-DECISION → JD-10 |
| 0008-0005/0006 | `reserved_tag_key` has no CLI code or exit mapping | gap | blocks-impl | joint | JOINT-DECISION → JD-8 |
| 0008-0005/0006 | `--tag recognized=` : programmer-mistake (0008) vs user-input (0005) | contradiction | blocks-impl | joint | JOINT-DECISION → JD-9 |
| 0008-0005/0006 | Reserved-key rule sits outside 0006's only blocking gate | gap | risks-impl | joint | JOINT-DECISION → JD-4 |
| 0009-0002 | Dump contract drops row kind + escape classes → **launders an escape-row breach** | round-trip | blocks-impl | joint | JOINT-DECISION → JD-11 |
| 0009-0002 | Normalized escape-row representation undetermined | gap | blocks-impl | joint | JOINT-DECISION → JD-11 |
| 0009-0002 | `writes = []` rejection is an undeclared 0002 failure mode | contradiction | blocks-impl | joint | JOINT-DECISION → JD-11 |
| 0009-0004 | 0009 files a binding obligation on 0004's implement stage; 0004 silent | gap | blocks-impl | joint | JOINT-DECISION → JD-7 |
| 0009-0004 | `<clear>` sentinel reaches 0004, which never mentions "clear" | contradiction | blocks-impl | joint | JOINT-DECISION → JD-12 |
| 0009-0004 | Unsanctioned escape writes pass read-back green | round-trip | risks-impl | joint | JOINT-DECISION → JD-12 |
| 0009-0005 | `GroupInternal` absent from 0005's declared mapping and code table | contradiction | blocks-impl | joint | JOINT-DECISION → JD-8 |
| 0009-0005 | New `CLIError` field vs 0005's "no new envelope fields" audit row | gap | blocks-impl | joint | JOINT-DECISION → JD-8 |
| 0009-0005 | Row-identity wire form has no declared inverse | round-trip | cosmetic | single-RDR (0009) | NO CONFLICT (cosmetic; blocks nothing) |

## Dispositions

### SPEC-DEFECT — RDR 0008 → Draft

**One** finding is genuinely internal to a single RDR rather than a shared
decision, and it is the one place a sibling's *verification* is false rather
than merely silent.

RDR 0008's assumption **A6** carries `Status: Verified`, `Method: Peer RDR`, and
this evidence:

> The string `Input` appears **zero** times in RDR 0009 (case-sensitive, whole
> file).

**A11** repeats it, adding "re-counted at Reconcile". Verified independently
during this pass:

```
grep -c '\bInput\b' 0009-escape-row-shape-conformance-ownership.md → 5
(lines 509, 736, 741, 742, 857 — including 0009's own Normative Contracts,
 where it reaches its subject through `in.Table`)
```

0008's citation `0009…md:536-540`, quoted as 0009's "at `Resolve` entry" text,
resolves in Final 0009 to unrelated prose about uSCXML conformance. The count
was true of the **Draft** 0009 that 0008 read; 0009 then gained A10 (from its
repeatability lens) which added exactly the `Input` references that falsify it,
and locked.

This is a SPEC-DEFECT, not a JOINT-DECISION: the wrongness is entirely internal
to 0008 (a false premise under a `Verified` stamp plus stale line citations),
0009 requires no change, and calling it "joint" would dodge a real demotion.
0008 is also the less foundational of the pair on this axis — 0009's text is
correct as written.

- **RDR demoted**: 0008
- **Target re-entry stage**: 4 (Resolve) — re-verify the two assumptions against
  Final 0009 and repair the citations.
- **SCOPE**: **STAGE-SCOPED**. The `re-verify` set is bounded (A6, A11). The
  chosen approach — reserving the key name and locating enforcement at the
  kernel input boundary — is untouched; what fails is the evidence for *where
  enforcement must live*, which A6 concluded from the false negative. A
  re-resolve that re-reads Final 0009 either re-derives the same locus on sound
  evidence or surfaces the overlap now recorded as JD-5. FULL-FLOW is not
  earned: no foundational assumption is voided and the Problem Statement stands.

### JOINT-DECISIONS — 12 tolerances, no further demotions

Every other blocking finding is a shared decision neither RDR solely owns:
each sits between a new producer and a peer that is *silent*, and silence is not
a defect the silent RDR can be charged with. Per Stage 7.1, these are marked and
homed, not absorbed as edits to a Final RDR's design body.

**The home for all twelve is a new umbrella RDR** (see *Required next action*).
There is no existing umbrella-decision record in this consumer, and 12
tolerances cannot be homed in a report that is itself only evidence.

| ID | Joint decision needing one normative home | Pairs |
| --- | --- | --- |
| JD-1 | Guard atom transport: how N atoms cross a 1-slot `Row.Guard`, and whether 0003's placement MUST admits `exists = false` in `all` | 0007-0003 |
| JD-2 | Resolver flow: gate-then-count vs match-then-escape; whether escape rules may carry guard blocks | 0007-0002 |
| JD-3 | `RequiresOwned` producer, preservation across the dump, and its value on escape rows | 0007-0002, 0007-0009 |
| JD-4 | Static/runtime correspondence: whether lint must detect unevaluable-at-runtime risk, incl. the absent `recognized` key and `contains`-over-absence | 0007-0008, 0007-0006, 0007-0003, 0008-0005/0006 |
| JD-5 | Precedence between the two whole-table `Resolve` entry preconditions (0009 breach, 0008 reserved-key) | 0007-0009, 0008-0009 |
| JD-6 | Which out-of-band channel a producer defect uses — typed error vs panic — and its reconciliation with 0005's never-silent envelope | 0007-0009 |
| JD-7 | Read completeness at the accessor→kernel boundary (0007 A6b): who states that a partial read takes the refusal branch | 0007-0004, 0009-0004 |
| JD-8 | CLI surface for the new refusals: `owned_state_unavailable`, `reserved_tag_key`, `GroupInternal`, and whether the envelope may gain fields | 0007-0005, 0009-0005, 0008-0005/0006 |
| JD-9 | `--tag` provenance: whether caller-supplied tags may satisfy owned/recognized state | 0007-0005, 0008-0005/0006 |
| JD-10 | Recognized-tag totality (absent vs declared-empty), the name constraint's effect on 0002's canonical fixtures, and fixture-rename ownership | 0008-0002, 0008-0009 |
| JD-11 | Escape-row identity across the dump: row kind + escape classes in the preservation set; `writes = []` | 0009-0002 |
| JD-12 | `<clear>` sentinel semantics at the write accessor | 0009-0004 |

## Verdict

**NOT RECONCILED.**

- One SPEC-DEFECT: **RDR 0008 → Draft**, re-entry at **Stage 4**, scope
  **STAGE-SCOPED**.
- Twelve JOINT-DECISIONS require a named normative home before the cluster
  implements. They do not demote anyone, but they are not yet *tolerances*
  either: a tolerance requires a home, and none exists yet. Recording them here
  is the mark; the hoist is the required next action below.

Do not implement over these. The concentration is diagnostic: 24 `blocks-impl`
findings across 13 pairs, every one at a seam where a new producer meets a
silent peer.

## Required next action (blocks the cluster)

1. **`/rdr-resolve 0008`** — re-verify A6/A11 against Final 0009, repair the
   stale citations, then `/rdr-finalize 0008`.
2. **Seed one umbrella RDR** to be the normative home for JD-1…JD-12, then add
   the `Final [joint decision → <umbrella §-anchor>]` qualifier to each affected
   sibling. Until that RDR exists, the twelve decisions are unhomed and the
   cluster's Status qualifiers cannot be written truthfully.
3. **Re-run this gate** after both, since the umbrella will touch cross-RDR
   seams.

## Observations outside this gate's authority

Recorded, not dispositioned — these are internal to single RDRs or to the
process, and no cross-RDR pair owns them.

- **RDR 0004 is `Final` with every Prerequisite unchecked**, including
  `- [ ] All Critical Assumptions verified`. Verified this pass. That is a
  finalize-gate question for 0004, not a cross-RDR contradiction.
- **RDR 0003's and 0002's Prerequisites are checked** and predate the
  obligations 0007/0008/0009 route to them; nothing invalidates a checked
  prerequisite when a later RDR binds it.
- **The whole-set critique** (`critique-set.md`, 34 rows) argues the set has
  substituted specification for delivery: 29,018 lines of RDR/evidence against
  3,971 lines of Go, with `version` as the only shipped verb, and three
  foundational RDRs seeded→finalized in one day. Its AT-1 ("a routed obligation
  must be accepted by its destination") is the generalized form of this
  report's set-level finding and is worth adopting as a finalize-gate item.
- **The prior `0001-0006` pass returned all-clear** on pairs this pass reopens
  (0007-0002 and 0009-0002 both land on 0002's resolver flow and dump contract).
  That gate was correct for the set as it stood; the drift arrived with the
  later locks — which is the case Stage 7.1 exists for.

## Review Gate

- **Cluster correct?** Yes. All eight members Final and unimplemented; 0001
  excluded as Implemented. Membership was widened mid-pass after a reference
  count showed 0007 citing 0004/0005/0006 thirty times, adding four pairs
  (`0007-0004`, `0007-0005`, `0007-0006`, `0008-0005-0006`) — three of which
  returned blocking findings that a narrower cluster would have missed.
- **No SPEC-DEFECT "fixed" by editing a Final RDR's design body?** Correct. The
  only edit made is 0008's status flip + README row + re-entry note.
- **Was the demoted RDR actually flipped and annotated?** Yes — 0008 carries
  `Draft [revised from Final 2026-08-11; re-verify A6, A11 — …]`, its README row
  is updated, and its re-entry note carries defect, target stage, scope, and
  direction inline.
- **Right RDR demoted?** Yes. The false premise and stale citations are 0008's;
  0009's text is correct as written and needs no change.
- **Each JOINT-DECISION genuinely joint?** Each names a decision neither RDR
  solely owns — in every case the peer is *silent* rather than wrong, so a
  demotion would punish silence. The one finding where a sibling's claim is
  affirmatively false (0008 A6/A11) is dispositioned SPEC-DEFECT, not joint.
- **Does each tolerance have a home?** **Not yet** — this is why the verdict is
  NOT RECONCILED rather than RECONCILED WITH TOLERANCES. The homes are specified
  as required next action 2.
- **SCOPE narrowest that covers the blast radius?** Yes. STAGE-SCOPED matches a
  bounded `re-verify` set (A6, A11) with the approach intact.
- **Joint-check line present on every member?** 0007/0008/0009 each carry
  `Joint-check: clear (7 peers)`. 0002-0006 predate that gate item; this pass
  ran the check for them against the cluster, and its result is the set-level
  finding above — the check they never ran is exactly the one that fails.
- **Both prompts run?** Yes — `critique-set.md` (whole-set) and 13
  `pairwise-<A>-<B>.md` files, all under this directory.
