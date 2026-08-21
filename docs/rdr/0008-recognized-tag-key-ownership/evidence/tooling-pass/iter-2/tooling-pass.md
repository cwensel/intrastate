Model: claude-opus-5[1m]

# Tooling Pass — RDR 0008, iteration 2 (post-demotion re-lock sweep)

Run as the mechanical pre-step of the Finalization Gate, after the 08.1 cluster
demotion, the Stage 4 scoped re-verify of {A6, A11}, and the Stage 6 iteration-2
reconcile. Regression focus: the re-verify and reconcile rewrote assumption
evidence, re-anchored 15 peer-RDR citations, and edited five body passages.

## CHECK 1 — Template section coverage

Every **Required** (spine) section is Present-substantive. Section-by-section
against `TEMPLATE.md`:

| Template section | Class | RDR |
| --- | --- | --- |
| Metadata | Required | Present-substantive (L11) |
| Problem Statement | Required | Present-substantive (L95) |
| Context / Background / Technical Environment | Required | Present-substantive (L169/171/201) |
| Research Findings / Investigation / Key Discoveries | Required | Present-substantive (L209/211/229) |
| Critical Assumptions | Required | Present-substantive (L293) |
| Proposed Solution / Approach / Technical Design | Required | Present-substantive (L1063/1065/1146) |
| Capability Dependencies | Required | Present-substantive (L1675) |
| Existing Infrastructure Audit | Required | Present-substantive (L1685) |
| Decision Rationale | Required | Present-substantive (L1693) |
| Alternatives Considered / Alt 1 / Alt 2 / Briefly Rejected | Required | Present-substantive (L1790–1867) |
| Trade-offs / Consequences / Risks / Failure Modes | Required | Present-substantive (L1876–2012) |
| Implementation Plan / Prerequisites / MVV / Phases | Required | Present-substantive (L2088–2188) |
| Validation / Testing Strategy | Required | Present-substantive (L2209/2211) |
| Finalization Gate | Required | Pointer line to `artifacts/gate.md` (the record, not a hollow section) |
| References | Required | Present-substantive (L2470) |
| Day 2 Operations | Conditional | Cleanly deleted — this RDR creates no persistent resource. PASS, not Missing. |
| New Dependencies | Conditional | Cleanly deleted — no dependency added or updated. PASS. |
| Performance Expectations | Conditional | Cleanly deleted — alternatives are not compared on empirical performance. PASS. |

`Critical Assumptions` sits as an `###` under Research Findings rather than a
top-level `##`. This is the settled corpus nesting, not a deviation — peer Final
RDR 0009 carries the identical structure (`0009…md:315`). Not a finding.

Placeholder / template-residue grep over the whole file: **0 hits** for
`_Draft placeholder._`, `seed skeleton`, `TBD`, `see above`, `[Conditional — `,
`[Resource]`, `[Capability]`, and bare bracketed template instructions. No block
from an older TEMPLATE.md survives.

**C1: no findings.**

## CHECK 2 — Method-label vocabulary

All 13 Evidence Records checked against the eight sanctioned labels. Compound
labels are members of the set joined with `+`, not paraphrases:

`Peer RDR` (L299) · `Source Search` (L364) · `Prior Art` (L411) ·
`Source Search` (L445) · `Derivation` (L512) · `Peer RDR + Source Search` (L572) ·
`Source Search + Peer RDR` (L658) · `Source Search + Peer RDR` (L721) ·
`Spike` (L777) · `Source Search + Peer RDR` (L830) ·
`Peer RDR + Source Search` (L903) · `Peer RDR + Source Search` (L979) ·
`Source Search` (L1032).

No missing, paraphrased, or off-vocabulary label. Records relabeled during the
re-verify (A6, A11) both remain in-vocabulary.

**C2: no findings.**

## CHECK 3 — Source Search self-reference

Every Source Search Evidence path resolved. Grep for
`0008-recognized-tag-key-ownership.md:` anywhere in the body: **0 hits**. No
Evidence resolves to this RDR or to its own artifact directory. The re-verify
re-pointed A6/A11 at `internal/resolve/resolve.go::Symbol` anchors and at RDR
0009 section headings — neither is a self-reference.

**C3: no findings.**

## CHECK 4 — Docs Only on load-bearing claims

Grep for `Docs Only`: **0 hits**. No record uses the label at all, so none can
lack a Spike or Source Search plan.

**C4: no findings.**

## CHECK 5 — Symbol resolution of anchors

Source anchors are `path::Symbol` as the doctrine requires, and all resolve in
the cited file at HEAD:

- `internal/resolve/resolve.go::Resolve` — resolves; `func Resolve(in Input) (Result, error)`.
- `internal/resolve/resolve.go::Input` — resolves; carries `Table`, `Flow`, `Owned`, `Observed`, `Recognized`, `Guards`.
- `internal/resolve/resolve.go::Row` — resolves; carries `RequiresOwned`, `Escape`, `Writes`.
- `internal/resolve/resolve.go::recognizedTagKey` — resolves; `const recognizedTagKey = "recognized"`.
- `internal/resolve/resolve.go::assemble` — resolves.

No phantom (resolves-nowhere) anchor. Peer-RDR anchors: the reconcile converted
all 15 live bare line citations to durable anchors (section heading or assumption
ID) — 0002 ×6, 0004 ×4, 0006 ×3, 0007 ×2. The two surviving bare-line 0009
citations sit **inside** the delete-on-re-lock defect record, where naming the
retired ranges is the record's purpose; they leave the file with that block at
re-lock. Line numbers alone are never checked and never block.

**C5: no findings.**

## CHECK 6 — Status consistency

Thirteen assumptions, **zero** `Pending` or `Unverified` — grep returns nothing,
so no settled-fact prose can lean on an unsettled assumption. Twelve carry
`Verified` (four with a scoping qualifier that narrows a consequence, not the
assumption); A12 carries `Refuted as stated` with its repair written as an
explicit hard Prerequisite gate item, which is a terminal disposition, not an
open one.

No checklist-vs-gate disagreement: the Finalization Gate section holds only the
pointer line, so there is no second copy of any assumption status to contradict
the first.

**C6: no findings.**

## CHECK 9 — Evidence-field budget (ADVISORY — never blocks)

One field over the 30-line budget, marginally:

- **A10** (L831) — 31 lines.

At the boundary: A8/A11 (L722, L582) sit exactly at 30. `1 field over budget,
1 line.` The mass is verification content — the accessor-half sweep and the
JD-9 scope correction the reconcile landed — and its load-bearing anchor
(`internal/resolve/resolve.go::assemble`) is findable in the first lines. The
prose has not outgrown the record. No action; reported per the check's contract.

## Verdict

**PASS** — no findings across C1–C6. C9 reports one field 1 line over an
advisory budget, which by the check's own contract does not affect the verdict.
Proceed to the Finalization Gate's written responses.
