# Iteration-2 brief — cluster 0002-0009 (Stage 7.1 re-run)

Model: claude-fable-5
Date: 2026-08-23

Read-only against the RDRs. Sub-agents write ONLY their own output file under
this directory. Never edit an RDR, the JDR, or the README.

## Paths (absolute)

- Records dir: `./docs/rdr/`
  - 0002 `0002-transition-table-as-reviewable-data.md`
  - 0003 `0003-guard-predicate-exhaustiveness.md`
  - 0004 `0004-accessor-execution-safety-model.md`
  - 0005 `0005-skill-integration-cli-contract.md`
  - 0006 `0006-graph-lint-authority-and-guarantees.md`
  - 0007 `0007-guard-predicate-totality.md`
  - 0008 `0008-recognized-tag-key-ownership.md`
  - 0009 `0009-escape-row-shape-conformance-ownership.md`
  - (ignore `*-postmortem.md`; 0001 is Implemented — code is truth, out of cluster)
- **Home** (joint-decision registry): `./docs/jdr/0001-resolve-kernel-seam.md`
- **Artifact of record for peer status**: `./docs/rdr/README.md` (Index table). All eight members are `Final` today.
- **Iteration-1 ledger**: `./docs/rdr/cluster-reconcile/0002-0009/report.md` (2026-08-11) — its findings table is the origin ledger. Prior per-pair scan files sit beside it (`pairwise-<A>-<B>.md`).
- Related later pass (0003·0006·0007 only, 2026-08-22): `./docs/rdr/cluster-reconcile/0003-0006-0007/reconcile-report.md`
- Pairwise prompt: `./.claude/skills/rdr-cluster-reconcile/pairwise.md`
- Whole-set critique prompt: `./.claude/skills/rdr-cluster-reconcile/2-critique.md` (whole-set variant)

## What changed since iteration 1

Every member moved (all eight have commits after 2026-08-11), and the home was
created and then answered several entries. Current home state (read it, do not
trust this summary alone):

- §D1 resolved → `Row` carries parsed atoms (closes JD-1, JD-6)
- §D2 resolved → gate-then-count (closes JD-2)
- §D3 resolved → read accessor returns complete tag set or refuses (closes JD-7)
- §D4 resolved → kernel enforces guard domain; per-atom evaluator seam; per-atom refusal payload replaces `Refusal.Guard` (closes JD-12)
- §JD-4 CLOSED 2026-08-22 (0003 records narrowing; 0006 cites; reuse `graph-unprovable-coverage`; 0006 finding record gains atom field)
- §JD-13, §JD-14 decided 2026-08-22 (single-valued field lives in 0003's declaration model; escape rows are ordinary participants in coverage union + overlap check)
- §JD-11 WITHDRAWN → routed as a single-RDR defect on 0002 (dump must preserve row kind + escape classes)
- STILL OPEN at the home: JD-3 (`RequiresOwned` producer), JD-5 (precondition precedence 0008 vs 0009), JD-8 (refusal codes; exit-3 call), JD-9 (`--tag` provenance), JD-10 (recognized-tag totality; fixture rename)

Status qualifiers today: 0005 `Final [joint decision → §JD-8, §JD-9]`; 0008 and
0009 `Final [joint decision → §JD-5]`; 0002/0003/0004/0006/0007 plain `Final`
(with lock notes), having re-walked and absorbed decisions.

## Per-pair task (each sub-agent)

For EACH pair assigned:

1. Run the pairwise prompt from `pairwise.md` verbatim against the two current
   RDR files. Anchor every finding to two direct quotes (or one + `silent`).
   Zero anchored findings is a valid result.
2. **Ledger trace.** Pull the iteration-1 rows for this pair from `report.md`.
   For each row, state one of:
   - `RESOLVED` — the home answered it AND each member's fenced (```normative)
     text is consistent with the answer (quote the consistent clause per member).
   - `ANSWER-CONTRADICTS` — the home answered it but a member's fenced text
     contradicts the answer (quote the contradicting clause; name the member).
   - `STILL-OPEN` — the home has not answered; the tolerance stands (say whether
     the sibling's Status carries the qualifier naming that question).
   - `SUPERSEDED` — the finding no longer exists in the current text (quote what replaced it).
   Also mark each NEW finding from step 1 as `LEDGER:<row>` if it traces to an
   iteration-1 row, else `NET-NEW`.
3. **Peer-status sweep.** Every sentence in either RDR asserting the *other*
   RDR's status (Draft/Final/Implemented) or a peer's line-number citation: check
   against README index / current text. A false one is a candidate CITATION
   REPAIR (cosmetic, single-RDR) — list them with line numbers.
4. **Restatement check.** Does either RDR restate (rather than cite) the
   mechanism prose of a homed decision (§D1–§D4, §JD-4, §JD-13, §JD-14)? Quote
   any restatement with line numbers.

Write `pairwise-<A>-<B>.md` (A = lower number unless the brief names the file
otherwise) into
`./docs/rdr/cluster-reconcile/0002-0009/iter-2/`,
starting with `Model: claude-fable-5`, then sections: Findings (prompt format,
each with `LEDGER:`/`NET-NEW` tag), Ledger trace table, Peer-status sweep,
Restatement check. Use `file:line` anchors throughout.

Return to the parent EXACTLY the rdr-common §return-packet:

```
verdict: PASS | BLOCK | INCOMPLETE | NEEDS_DECISION
blocking: yes | no
evidence_paths: [...]
changed_paths: [...]
next_action: ...
summary_50w: ...
```

followed by a compact findings table (one line per finding:
`pair | TYPE | SEVERITY | OWNERSHIP | LEDGER/NET-NEW | 12-word gist`) and the
ledger-trace lines (`row | RESOLVED/ANSWER-CONTRADICTS/STILL-OPEN/SUPERSEDED`).
PASS = no blocks-impl/risks-impl single-RDR finding and no ANSWER-CONTRADICTS;
BLOCK otherwise.
