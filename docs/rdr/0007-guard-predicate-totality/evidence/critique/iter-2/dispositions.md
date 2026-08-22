Model: claude-opus-5[1m]

# Critique resolve — disposition ledger (iteration 2)

Origin ledger = the merged queue in `diff-modelB.md`, built by passage anchor
over `critique.md` (pass A, `claude-opus-5`, C-1..C-14) and `critique-modelB.md`
(pass B, `claude-fable-5`, B-1..B-13). Dual-model obligation satisfied at
generation; the diff ran in a third context that authored neither pass.

**Zero directional disagreements** between the passes — every shared passage got
the same verdict, differing only in lumping. Convergence on an independent draw
is treated as a hotspot signal, not redundancy.

Grounding gate run against: code on the working branch `via-claude`
(`internal/resolve/` is absent from `main`, which is docs-only), the peer RDRs
and JDR 0001 under `docs/`, and the RDR's own decided text.

## Dispositions

| # | Queue item (A/B ids) | Disposition | Section touched |
|---|---|---|---|
| 1 | Existence-polarity inversion in the handed-forward spike probe (C-1, B-2) | **fixed** | A3 Coverage gap — probe scoped to conjunction cases only; existence case named wrong and its rewrite specified |
| 2 | A3 spike proved a superset, not the specified shape (C-6, B-2) | **dismissed-with-cite** | — already stated: A3 Status, Testing Strategy preamble, and Phase 1's exit condition all name this; cove iter-2 F-2 flipped A3 to Pending for it |
| 3 | Resolution-scope veto has no writable mitigation (C-9, B-5) | **dismissed-with-cite** | — decided text: Load-Bearing Decisions rejects the narrowed veto; Failure Modes states the `unless` footgun and its idiom; A12 carries the open half. Re-raise of a settled call |
| 4 | Payload CLI transport undecided at lock (C-5, B-3) | **dismissed-with-cite** | — A19 already states the reversal, the fallback, and the reversibility argument; Prerequisites already scopes the block to 0005's rendering |
| 5 | A21 claims the kernel enforces no part of 0009's escape shape (C-3, C-4) | **fixed** | A21 Evidence corrected against 0009's kernel-precondition block; new *Fixture collision* note; A25 added (Pending); Testing row 15 rewritten; Prerequisites 0009 item extended |
| 6 | Flagship refusal's winning class has no renderable Code (C-10, B-4) | **dismissed-with-cite** | — the `disposition` mini-check already carries both exit-group notes and routes them to 0005's re-lock |
| 7 | Producer-half guarantees ride on unlanded 0002 duties (C-8, B-1) | **dismissed-with-cite** | — A16/A22/A24 each state their producer half as Pending; Prerequisites carries the two 0002 duties and A22's fallback |
| 8 | Five external prerequisites, several unsequenced (C-14, B-11) | **fixed** (partial) | Prerequisites — 0009 item now sequences the fixture collision; new A9 item added. The count itself is dismissed: each box is individually justified and A19's non-blocking status is argued |
| 9 | Nil-seam widening has no discriminating test (C-7) | **fixed** | Testing Strategy row 19 — leg marked REQUIRED, both assertions specified, silent-plan-production risk named |
| 10 | Empty `unless` algebra conflicts with 0003's subtraction (C-11) | **fixed** | A9 Status downgraded; Evidence corrected (inference, not citation); *Downgraded* note added; Prerequisites gains the 0003 item |
| 11 | New exported surface unswept against frozen boundary tests (C-12) | **fixed** | A26 added (Pending) naming the three banned lists and the REQ-37 constraint on A22's fallback |
| 12 | SQL:2003 supports the atom rule, not the veto (C-13) | **fixed** | Aggregation clause — prior-art scope stated explicitly; SQL named as the resolution-scope counter-example |
| 13 | `--tag` bypass unauditable (B-7) | **dismissed-with-cite** | — A13 is an accepted exposure with home §JD-9; Failure Modes states it; Testing row 16 pins the disposition |
| 14 | Masking path reopens via `Match` tags (B-8) | **fixed** | Failure Modes — new entry; grounded on `TagSet.matches` `!ok → false` |
| 15 | Closed reason set can't distinguish who-can-fix-it (B-9) | **dismissed-with-cite** | — Failure Modes already states the skew case and why the reason field bounds it; widening the closed set is net-new scope against a set this RDR fixes by name |
| 16 | Pruned rows report nothing (B-12) | **fixed** | Failure Modes — new entry; states the accepted cost and the anti-pattern pressure |
| 17 | `contains` present-key contract unwritable at lock (B-13) | **dismissed-with-cite** | — Phase 3 already states the block, its cause, and what is unaffected; Phase 4 carries the declaration request |
| — | Payload sort-key totality over set literals (C-2) | **fixed** | Payload clause — two stated conditions added (literal byte-spelling; dedupe is an emit obligation, not a sort property) |

Fixed 8 · dismissed-with-cite 9 · charted 0.

## Charted

Nothing charted this pass. Every finding either landed inside the current
design's scope or was already adjudicated in the draft; no finding required
expanding the RDR's remit.

## Needs (re)verification — carried to Stage 6

New or disturbed load-bearing claims. None verified here by charter.

- **A25** (new, Pending) — RDR 0009's kernel entry precondition vs. Phase 1's
  fixture migration. The collision is confirmed on source; the SEQUENCING is a
  cluster decision this RDR does not settle.
- **A26** (new, Pending) — the exported-name sweep against REQ-25/36/37. Folds
  into A3's re-spike. Also records that REQ-37 forecloses a kernel-side
  `Canonicalize*` name, which constrains A22's fallback.
- **A9** (downgraded) — Verified for the kernel verdict rule; the agreement with
  RDR 0003's subtractive lint algebra is Pending and rides to 0003's re-lock.
- **A3** (scope sharpened, still Pending) — the re-spike must additionally
  rewrite the probe's existence case to carry an explicit literal and assert
  both polarities.
- **Payload sort totality** — now conditional on RDR 0003 declaring a canonical
  literal spelling for set-valued literals (`in`, `contains`). Same 0002/0003
  canonicalization duty as A22/A24, applied to literals.

## Mini-checks

Cue read re-run against this pass's fixes. The four cues that fired on this
draft (`authority`, `oracle`, `disposition`, `trace`) already carry their tables;
`fidelity` still does not fire (no inverse invariant is claimed). **No fix in
this pass added a new cue** — the edits corrected claims, sharpened test
obligations, and added assumptions, none of which introduces a fallback path, a
new source-of-truth candidate, an absence-of-error oracle, or a new
output-surface assertion pair. Tables persist in the draft unchanged.

One `authority` row is now stale in the reader's favour rather than against it:
the *Tag-key canonicalization* row's "duty not yet a 0002 clause" now also covers
literal spelling (payload clause). Not re-tabled — the row's writer/canonical
columns are unchanged.

## Review gate

- Findings anchored to real passages, named symbols, and quoted clauses across
  both passes; origins span `§1`/`§2`/`§3`/premortem/AT rows. Healthy — no
  re-run on another model owed.
- Fixes are grounded: every `fixed` row cites source confirmed this pass
  (RDR 0009's kernel block, `escapeRow`, the spike probe's assertion,
  `TagSet.matches`, `exportedKernelSymbols`' banned lists, RDR 0003's operator
  matrix and subtraction sentence, RDR 0005's stable-code table).
- Net-new scope was not folded in: items 3, 13 and 15 would each have widened a
  contract this RDR fixes by name, and were dismissed against decided text.
- No tiebreaker surfaced — every fork collapsed on the evidence.
