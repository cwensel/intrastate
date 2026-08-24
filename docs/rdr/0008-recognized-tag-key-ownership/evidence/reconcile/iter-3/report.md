Model: claude-opus-5[1m]

# Stage 6 Reconcile — RDR 0008, iteration 3 (second post-demotion re-entry)

Verdict: **RECONCILED**. Thirteen assumptions terminal, no BLOCKER, no unrun
spike. This pass follows the `0002-0009` cluster iteration-2 demotion and the
scoped Stage 4 re-verify of {A9, A4} against Final RDR 0002. Iterations 1–2
dispositions are undisturbed except where the rewrite touched them.

## Stage 5 preflight

- `Profile: foundational` (recounted at the re-entry: six `normative` blocks,
  all clauses of one reserved-key identity rule) → required lenses
  `cove 3amigo critique repeatability`. All four evidence dirs present and
  complete.
- Determinacy / repeatability: `full` variant as `foundational` requires — all
  three runs stamp `variant: full (profile: foundational)` with three
  **distinct** models (`claude-opus-5[1m]`, `claude-fable-5`, `glm-5.2:cloud`),
  plus `diff.md` and `resolve.md`. No variant mismatch.
- **Caveat carried from iterations 1–2** (not a gate failure): `critique` ran
  single-model (Pass A + Pass B, both `claude-opus-5[1m]`, fresh contexts) as
  the fallback `resources.md` records; `iter-2/delta-verdicts.md` likewise.
  Repeatability supplied the genuine multi-model draw.
- The demotion invalidated no lens: the re-entry was declared STAGE-SCOPED to
  {A9, A4}, and no lens finding rested on the withdrawn pre-lift reading.

## Open set — four sources

1. **Pre-Lock lists**: fully dispositioned at iterations 1–2; nothing reopened
   by this demotion.
2. **Still-Pending assumptions**: none. All thirteen carry a terminal Status
   (grep for `Pending`/`Unverified`: 0).
3. **Named-but-unrun spikes**: none. Two spikes exist and both have captured
   output — `spikes/a9-normalization.md` (iteration 1, now superseded) and
   `spikes/a9-lifted-outcome.md` (this re-entry, five runs with verbatim
   output).
4. **Exactness-word delta** — the post-mutation set introduced by the 413-line
   A9/A4 rewrite. Each claim traces to a quoted Final-0002 fence, a captured
   spike run, or a rule this RDR authors; the load-bearing ones were
   independently re-derived from source at this stage (below).

Absorption audit over all six evidence dirs, scoped to lens findings whose
landing site fell inside a rewritten passage: **one partial regression**
(critique D-17), fixed here. Everything else stands — cove F-4/F-5/F-7/F-8/F-9,
3amigo IMP-1/IMP-5/QA-2/QA-3/PM-1/PM-3/PM-4, critique D-1/D-2/D-9/D-11/D-13/
D-22 and the iter-2 block-4 interim rule, repeatability D-2. Blocks 3–6 were
untouched by the rewrite and are intact.

## Dispositions

| Item | Source | Disposition | Evidence / plan |
|---|---|---|---|
| **A9** — what Final 0002's normalizer does with the reserved key | 1,3,4 | **VERIFIED (restated; prior reading withdrawn)** | Every fence A9 quotes is verbatim in Final 0002 — outcome binding "reads the match blocks only" / "MUST bind exactly one outcome"; guard-position atom "MUST be refused at load … never lifted"; "Load is fail-fast: the first category a document trips is the refusal". 0002 Status confirmed `Final`. Spike `spikes/a9-lifted-outcome.md` §1 reproduces 0002's `output.txt` byte-for-byte. Kernel side re-confirmed: `internal/resolve/resolve.go::recognizedTagKey`, `::assemble` injection. |
| **A4** — census + migration inventory | 1,3,4 | **VERIFIED (census corrected; inventory empty)** | Re-derived from disk at this stage, not taken from the re-verify: 0002's `rdr-fixture.toml` and `kata-fixture.toml` both carry `[tags.recognized]` with `provenance = "recognized"` and 4+2 `[rule.match.recognized]` sites — the JD-10 rename is landed. The single remaining `recognized`-provenance guard atom is `0003…/guard-fixture.toml:36-37,72`, refused at load (spike §5), and is RDR 0003's artifact. Second clause: the Go literal occurs once as a read-only match pattern (`fixtures_test.go:237`) plus the kernel const. |
| **A9 observation** — 0002's spike lifts from any block | 3 | **ACCEPTED (Design Decision — booked on the peer)** | Not routed and not a defect record: 0002's own Testing Strategy scenario 3 already names both guard-position mutants among those owed and records the spike as a nine-category witness, remainder owed at implementation. 0008 restates 0002's fenced rule, not its spike's current behavior. Editing a `Final` peer from a sibling's stage is forbidden regardless. |
| **A10** — accessor-derived key spelling | 2 (carried forward) | **VERIFIED (re-checked against the rebuilt peer artifact)** | A10's evidence cites the *same* 0002 spike the rebuild moved, so it got the "peer is moving" check the ledger says was missing. At HEAD `main.go::Accessor` is still `{Mode, Path}` — no key field — and the tag→accessor reference is still one-way (`tag.Accessor` names an accessor). Claim survives the rebuild unchanged. |
| **A1** — additive extensibility of 0002's category list | 2 (carried forward) | **VERIFIED (re-checked)** | The "including at minimum" fence is verbatim in Final 0002; `reserved_tag_key` now appears there in the category block, the ownership table, and the fenced text. A1's Overrides-licensed half is ratified by §JD-10. |
| **A12** — handoff to 0002 | 1,2 | **DOWNGRADED → discharged in the carried form** | Stronger than the Prerequisite claims: Final 0002 names `reserved_tag_key (RDR 0008)` in its validation-category block, carries the naming rule in its own fenced text, and lists "Reserved recognized tag key \| RDR 0008" in its dependency table. Residual — block 3's payload fields + advisory carrier — is unabsorbed by 0002 and rides §JD-8 as an open Prerequisite gate item, correctly unchecked. |
| **JD latches** (JD-4, JD-5, JD-8, JD-9, JD-10) | 4 | **VERIFIED** | Re-read at the JDR: JD-10 **ANSWERED 2026-08-23** (totality ratified from 0002's re-lock); JD-4 **CLOSED 2026-08-22**, and the RDR's "answered by §D4" citation resolves (`§D4` exists). JD-5, JD-8 (widened 2026-08-23), JD-9 open and cited as open. Status-line qualifier `[joint decision → §JD-5, §JD-8]` matches. |
| **Sweep-on-the-way items** named by the demotion's Direction | 4 | **VERIFIED (landed)** | "0009 still `Draft`" corrected — 0009 is `Final`; stale 0007 quotes re-anchored (0007 is `Final`, its A13 no longer carries the quoted sweep); "§JD-4 … open there" corrected to answered-by-§D4 at both sites. |
| **Block 2 / scenario 5 restatement** (the demotion's two contradictions) | 1,4 | **VERIFIED (restated to the peer's shape)** | Block 2 now defers every predicate position to 0002's outcome-binding contract, states the lifted atom and the guard-position load refusal, and attributes the lower bound to 0002 rather than claiming none. Scenario 5 now expects **exactly one** refusal by category with the unspecified-which-one clause, citing the fail-fast fence, and names spike §3 (`reserved_tag_key: recognized declaration named "outcome"`, exit 1) as its normative fixture. |
| **Phase 2 rename scope** | 4 | **VERIFIED (empty)** | No rename remains for this RDR: 0002's are landed, 0003's fixture is 0003's and is a load failure, not a rename. Author's round approved dropping it (3/3 approved, 0 rejected). |
| **critique D-17** — deciding-rows paragraph vs. the swept inventory | audit | **VERIFIED (repaired here)** | The one un-absorbed lens fix. D-17 exists so the QOC deciding-rows prose states R's *actually shipped* cost; the rewrite swept the rename to empty at the blast-radius table row, A4, Consequences, Risks and Phase 2 but missed this paragraph, which still read "a three-fixture rename (A4)" — contradicting the table row twelve lines above it. Restated to the JD-10 rename discharged at 0002's re-lock with an empty inventory. §amendment-sweep run: no other stale count survives (line ~144's "all three committed" is explicitly historical, "before RDR 0002's re-lock"). |
| A2, A3, A5, A6, A7, A8, A11, A13 | 2 | **VERIFIED (undisturbed)** | Out of the re-entry's scope and untouched by the rewrite; lock-time stamps carry forward per the STAGE-SCOPED declaration. A6/A11 were themselves re-derived at iteration 2. |

## Hard rules

- **Refutation → BLOCKER?** No. A9's *prior* reading was refuted, but that
  refutation is what the Stage 4 re-entry already processed: the assumption is
  re-verified in its restated form, and the RDR's approach is untouched — the
  name is still reserved and still enforced at the kernel input boundary. The
  refutation moved the RDR **toward** its peer, not away: Final 0002 has since
  adopted this RDR's naming rule and category into its own fenced text and
  names 0008 as their owner. Nothing routes back.
- **MVV-critical deferral?** None. The kernel half runs during this RDR's
  implementation. The normalizer half is scoped to the peer that owns the code
  and is gated by a written Prerequisite. The one open Prerequisite (§JD-8
  payload/advisory carrier) gates scenarios 4 and 7 on *0002's* side, not what
  the MVV proves.

## Flow hygiene

The §punt-ledger row for this route-back was appended by Stage 4 and is
present (third row). Its Notes block still read "Both rows share one root
cause" and "Neither escape was caught by a pre-lock lens" against three rows;
corrected here to three, with the third cause stated (0002's spike read as it
stood, never re-read after the rebuild and re-lock).

## Completeness

- No `_Draft placeholder._`, no seed-skeleton header (grep: 0).
- No surviving template brackets (`[Required`/`[Conditional`/`[Resource]`/
  `[Capability]`: 0).
- `## References` filled with real citations and durable anchors.
- Thirteen assumptions, zero `Pending`/`Unverified`.
- Remaining bracketed text is the Finalization Gate's own response prompts —
  Stage 7's to answer, correctly unfilled here.
- The Status qualifier and `## Refinement Context` correctly remain: both are
  cleared by the Stage 8 re-lock, not by Stage 4 or 6.

## Carried to Stage 7 (not blockers)

- Final 0002's dependency table records RDR 0008 as `Final` while 0008 is
  `Draft [revised from Final]`. It is a stale fact on a `Final` peer, not
  amendable from this RDR's stage, and it self-corrects at 0008's re-lock.

## RDR edits made by this stage

Decision Rationale's deciding-rows paragraph restated from "a three-fixture
rename (A4)" to the JD-10 rename discharged at 0002's re-lock with an empty
migration inventory (critique D-17 re-absorbed); postmortem ledger Notes
corrected from two rows to three with the third root cause stated.
