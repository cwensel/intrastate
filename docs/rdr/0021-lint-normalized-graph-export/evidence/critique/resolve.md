# Critique resolve — cli/0021 (iteration 1)

Dual-model: `critique.md` (claude-opus-5) + `critique-modelB.md` (claude-sonnet-5).
Ledger diff by element anchor (`C-N` ids are per-file; the element id is the join key).

Converged on both passes (hotspots): 0021:A8, 0021:C2, 0021:C4, 0021:C5.
Opus only: A1, A6, RT1, S6, §metadata, §pre-lock-mini-checks.
Sonnet only: A2, C1, §context, §decision-rationale, §failure-modes,
§problem-statement, §references.

## Dispositions

| Row | Anchor | Disposition | Basis |
| --- | --- | --- | --- |
| A C-2 / B C-1 | 0021:A8 | **fixed** | Founded. Prerequisites gated only A1/A2/A5/A6; A8 carries C4's MUST, S7 and the disposition-table ceiling row. Added an A8 Prerequisites gate on Phase 2. |
| A C-3 | 0021:C2 | **fixed** | Founded in code: `reach.go::OpaqueValue = "<opaque>"`, returned at `reach.go:438` into `Node.Values`, so it reaches `values[]`. C2 now names it; new assumption A9 books the collision question. |
| A C-7 / B C-5 | 0021:C5 | **fixed** | Founded. C5 promised "one `jq -r` step" but named no key. Now normative: `data.dot`, unwrap `jq -r .data.dot`. Swept S5 to agree. |
| A C-10 | 0021:RT1 | **fixed** | Founded. "value identity on the exported projection" was self-scoped; requantified over C2's field list. |
| A C-1, C-9 | 0021:C4, §metadata | **dismissed-with-cite** | Unfounded. C4 mints `graph-export-too-large`; `graph-product-too-large` is the guard-product code (`guard/lint.go:22`, `taxonomy.go:34`) on a distinct bound (`ProductBound()` vs `nodeCeiling=4096`). C4 already states the distinction. 0029 (Final) explicitly disclaims minting it — its census defers the tier to cli/0021. No locked record is touched. |
| A C-4 | 0021:A6 | **dismissed-with-cite** | Re-raise. A6's Evidence already records the spike spellings as provisional placeholders and states "C2's spellings govern." |
| A C-5 | 0021:S6 | **dismissed-with-cite** | Unfounded. S6 cites A6's recorded escaping ORDER (backslash→quote→newline), which A6 states is sound and carries into implementation — not a deleted /tmp artifact. |
| A C-6 | 0021:C2 | **dismissed-with-cite** | Re-raise of the Q2 ruling (rulings.md): `[]`/`{}` never `null`, optional member ABSENT. C2 states it; the `disposition` mini-check table pins both rows. |
| A C-8 | 0021:A1 | **dismissed-with-cite** | Refuted by A1: the spike byte-compared a multi-kilobyte multi-line payload through `TextLiner` (`respond.go:58-65,165`), advisories on stderr. |
| B C-2 | 0021:C1 | **dismissed-with-cite** | Not a defect. `--flow` resolving no ids this build is deliberate and stated in C1; mirroring lint's arm set verbatim is the contract. |
| B C-3 | 0021:C2 | **dismissed-with-cite** | Out of frame: authored-TOML fidelity is A3's subject, Verified by source search with no re-parse; C2 predicates carriage of the document. |
| B C-4 | 0021:C4 | **dismissed-with-cite** | Re-raise. C4's oracle is stated mechanism-independently and binds the in-traversal-observer arm too (P-14); A2 names that arm as the fallback. |
| B C-6 | 0021:C2 | **dismissed-with-cite** | C2 fixes exact field spellings as normative and S2 pins golden fixtures; a schema artifact is not owed at lock. |
| B C-7 | 0021:§decision-rationale | **dismissed-with-cite** | Re-raise of the settled surface choice (ALT1/ALT2 weighed; Q3(c) ruling). |
| B C-8 | 0021:C4 | **dismissed-with-cite** | Re-raise: the refusal-over-partial-document call is stated in C4 and accepted in F1 (premortem P-13) with the narrow-a-domain remedy. |
| B C-9 | 0021:§failure-modes | **dismissed-with-cite** | F1 already records the killed-process case: no terminal envelope, existing contract. |
| B C-10 | 0021:§problem-statement | **dismissed-with-cite** | The abstraction marker is REQUIRED in C2 and rendered into DOT; the soundness sentence states which claim classes survive. 0022 is a named non-blocking consumer. |
| — | `taxonomy.go:150-152` | **charted-to-successor** | Stale doc comment attributes the node ceiling to `graph-product-too-large` (see charted.md). Source defect, not an RDR defect. |

## Needs (re)verification — Stage 6

- **A9 (new, Pending)** — `<opaque>` reaches the wire as an ordinary `values`
  member; verify no authored value can collide with that spelling.
- **A8 (already Pending)** — now explicitly gates Phase 2 in Prerequisites.
- C5's `data.dot` key is a new normative wire name; MVV step 4 and S5 assert it.

## Tiebreakers

None. Every fork collapsed on evidence.
