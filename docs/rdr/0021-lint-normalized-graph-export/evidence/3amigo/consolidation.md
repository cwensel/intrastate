# 3amigo consolidation — 0021-lint-normalized-graph-export

Model: claude-opus-5 (consolidator; persona passes all `claude-sonnet-5`)
Date: 2026-09-12
Iteration: 1 (base)

Three isolated persona passes, no cross-persona visibility. Consolidation is
mechanical: hotspots are `rdr anchors --record 0021` over each persona file,
counted, `>=2`. Overlap marks a hotspot PASSAGE, not a validated finding; a
finding raised by exactly one persona is not thereby weaker.

## Isolation check

Anchor sets are near-disjoint along persona ownership (PM → prose sections;
implementer → `C*`/`D*`; QA → `S*`/`A*`/`MVV`). No persona file cites another
persona's output. Isolation held.

## Hotspots (cited by >=2 isolated personas)

| id | Personas | What converged |
| --- | --- | --- |
| `0021:C2` | implementer, QA | The `reach.nodes[].values` shape is unfixed (object vs array); both reached it independently, and the record's own `trace` mini-check already flags it as CONTRADICTION (open as A7). |
| `0021:C5` | implementer, QA | Mode-coexistence clause: implementer on the unresolved `JDR 0002 §D1` citation; QA on whether the JSON string member carries the trailing newline A1 pins for text mode. Two different defects on one passage. |
| `0021:MVV` | PM, QA | The acceptance floor does not validate the stated user outcomes — PM on diagram legibility, QA on goldens having no fixed shape to pin. |
| `0021:S6` | PM, QA | DOT hostile-content scenario: PM on it standing in for diagram-review, QA on its oracle being exit-code-only against a known escaping bug. |

## Merged ledger

Origin ledger for the resolve half. `P1`=PM, `P2`=implementer, `P3`=QA.

| # | id | Persona | Finding | Severity |
| --- | --- | --- | --- | --- |
| 1 | `0021:C2` | P2, P3 | `reach.nodes[].values` shape contradicted between Illustrative Code (object keyed by tag) and the A5 encoder fixture (array of assignment strings); blocks the Phase 1 struct definition and leaves S2/S6/S9/MVV goldens with no fixed shape. **HOTSPOT.** Tracked as A7 (Pending). | high |
| 2 | `0021:C4` | P2 | The incompleteness refusal (`graph-export-too-large`) depends on the `complete` bool from private `reach.go::reach`, but exported `graphlint.Reach` discards it — no exported symbol surfaces completeness. Contract names a refusal with no reachable trigger. | high |
| 3 | `0021:S6` | P3 | Oracle is "survives `dot -Tsvg`" (exit 0) + id-set match — cannot catch a well-formed-but-mis-escaped DOT label, and A6's own evidence names a real escaping bug carved out as "styling, not fixture." **HOTSPOT.** | high |
| 4 | `0021:A7` | P3 | Pending assumption leaves the `values` shape undetermined while C2/S2/S6/S9/MVV treat it as settled; the record's own trace table calls it a contradiction. (Same defect as #1, viewed from the assumption side.) | high |
| 5 | `0021:C5` | P2 | Forward-compat clause cites `JDR 0002 §D1`, unresolved in this record's edges and ambiguous against the RDR 0002 cited throughout. **HOTSPOT.** | medium |
| 6 | `0021:S5` | P3 | Claims byte-for-byte equality between text-mode DOT and the `jq -r`-unwrapped JSON DOT without pinning whether the JSON string member includes the trailing newline A1 pins for text mode. **HOTSPOT passage (C5).** | medium |
| 7 | `0021:C1` | P2 | Unclear whether "mirrors lint's arm set verbatim" extends to `--flow`'s reserved-flag help text or only to the error codes. | medium |
| 8 | `0021:S1` | P3 | Asserts no map-order dependence but names no provocation mechanism (e.g. `GODEBUG=randmapiter=1`); a small fixture can pass by luck rather than by proof. | medium |
| 9 | `0021:S8` | P3 | Refusal-arm scenario does not itself assert "no artifact on stdout"; that lives only in the mini-check `disposition` table, not in the scenario text. | medium |
| 10 | `0021:§problem-statement` | P1 | Motivates the RDR partly on handing the export "to an external formal checker," which `§background` puts out of scope and no test validates. | medium |
| 11 | `0021:MVV` | P1 | Diagram-review outcome validated only by DOT rendering + id-set equality (S6), not legibility/usability for a human reviewer. **HOTSPOT.** | low |
| 12 | `0021:§decision-rationale` | P1 | CI-diffability claim has no direct "diff across commits" scenario; determinism and the golden pin are proxies, not the outcome. | low |

Informational, not a defect: P2 verified `guard.AssignmentCount`'s signature
matches C2's finite-domain test as written (clean).

## Dispositions (resolve half, iteration 1)

Grounded against code on `main`, `{RDR_RESOURCES}`, and the RDR's own decided
text before any edit. `charted.md` holds the charted rows.

| # | id | Disposition | Basis |
| --- | --- | --- | --- |
| 1, 4 | `0021:C2`, `0021:A7` | **fixed** | `reach.go::Node` carries `Values map[string][]string`, so the OBJECT exhibit projects without invention; the array-of-`key=value`-strings form exists only inside `(Node).key`'s escaped fingerprint. C2 now pins the shape in both the `reach` description and the field-spelling list; A7 → Verified; the A5 fixture is named the wrong exhibit; the `trace` CONTRADICTION row now reads RESOLVED. |
| 2 | `0021:C4` | **fixed** | Confirmed: `reach()` returns `(nodes []Node, complete bool)` but `Reach` is `nodes, _ := reach(m)`, and no exported symbol in the package surfaces completeness. C4 now requires an exported completeness surface (the Q3(c) edge function), cites `graphlint.NodeCeiling()` instead of the private symbol, and distinguishes `reach`'s node-count completeness from `analysis.go`'s `checkNodeCeiling` guard-product bound. New assumption A8 (Pending) books the claim. |
| 3 | `0021:S6` | **fixed** | A6's own Evidence names a real escaping bug carved out as styling, and the oracle was exit-0 + id-set only. S6 now also asserts each hostile label DECODES back to its source tag value, so a well-formed-but-mis-escaped label fails. |
| 5 | `0021:C5` | **dismissed-with-cite** | `docs/jdr/0002-success-envelope-projection.md` §D1 "The caller-projection partition doctrine" exists and governs projection-before-`respond.OK`, echo-group-only, enforced partition, always-keep-core — exactly what C5 cites it for. Not RDR 0002; JDR is a separate namespace the RDR edge facet does not index, which is why it read as unresolved. Citation correct. |
| 6 | `0021:S5` | **fixed** | F1 pins text-mode stdout as `document + "\n"`; the enveloped string member carries no gateway newline, so "byte-for-byte" was off by one byte as written. S5 now compares against the document and names the newline explicitly. |
| 7 | `0021:C1` | **fixed** | C1 enumerates code spellings, so mirroring was always about codes; "verbatim" left help text ambiguous. Narrowed in C1 to the arm set and code spellings, with the `authority` table and the Existing Infrastructure Audit row swept to match (§amendment-sweep — three sites, one wording). |
| 8 | `0021:S1` | **fixed** | A2's spike did provoke map order under `GODEBUG=randmapiter=1`; only the scenario text omitted it. S1 now names the mechanism. |
| 9 | `0021:S8` | **fixed** | "No document on stdout" lived only in the mini-check `disposition` table. S8 now asserts it, so the Testing Strategy is self-sufficient. |
| 10 | `0021:§problem-statement` | **dismissed-with-cite** | Re-raise of decided text: `§background` scopes out "shipping or maintaining a TLA+/SMV model in-repo — a generator script over the JSON is a downstream consumer," and `§consequences` frames formal models as "derivable downstream consumers." The problem statement claims derivability, not a formal-checker feature. |
| 11 | `0021:MVV` / `0021:S6` | **charted-to-successor** | DOT styling is explicitly non-normative (Load-Bearing Decisions; F3 scopes the fixture to the node/edge SET), so a legibility criterion would reverse a decided call. See `charted.md`. |
| 12 | `0021:§decision-rationale` | **charted-to-successor** | A true cross-commit diff needs two builds, which C3's narrowing to (model, build) places outside this RDR. See `charted.md`. |

## Needs (re)verification — carried to Stage 6

- **A8 (new, Pending)** — the exported `graphlint` surface returning completeness
  alongside nodes and edges, addable without touching `reach()` or lint's call
  site. Raised by the C4 fix (a new load-bearing claim about a public surface).
- **A7 (Verified this pass)** — flipped Pending → Verified from source, not from a
  spike; the wire shape it fixes is normative in C2, so Stage 6 should confirm the
  A5 fixture was actually corrected rather than merely relabelled here.
- **Determinacy line** — written as `fired` in Normative Contracts this pass
  (was unjudged), which routes the repeatability add-on for a `large` profile.

## Loop verdict

`rdr-loop.toml --outcome lens-loop` with `iter=2 fix=substantial found=some
net_new=none` → **`rerun`** (`rule: loop-rerun`): the fixes rewrote normative
clause text in C2, C4, C1 and four scenarios, and a rewrite can open gaps.
Pass 2 is delta-scoped to the still-open anchors above, writing under
`iter-2/`. Not converged; the next lens does not start yet.

