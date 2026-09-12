Model: claude-sonnet-5

Persona 2 — Implementer review of RDR 0021 (docs/rdr/0021-lint-normalized-graph-export.md)

Widened beyond the owned C/D set to: the Pre-Lock Mini-Checks `trace` row (already self-flags A7), A2/A5/A7 (source-anchored assumptions C2/C4 depend on), Existing Infrastructure Audit, and source at internal/graphlint/reach.go, internal/cli/lint.go, internal/guard/declaration.go, internal/guard/product.go, internal/cli/clierr/clierr.go. Sent there by: C2's `reach.nodes[].values` shape is contradicted by the record's own two exhibits, and C4's incompleteness check cites a private function (`reach.go::reach`) that no current exported symbol surfaces — both are first-hour blockers a contract-only read cannot resolve.

---

Findings, severity-ranked:

1. [High] 0021:C2 — The wire shape of `reach.nodes[].values` is undetermined at lock, and the RDR admits it: the Illustrative Code renders `values` as an OBJECT keyed by tag (`{"stage":["draft"]}`), matching `graphlint.Node.Values map[string][]string` in reach.go, but the A5 encoder fixture renders it as an ARRAY of assignment strings (`["env=dev","flag&x"]`). C2's own field list says "`reach{… nodes[{id, values}]}`" without fixing the inner shape. This blocks Phase 1 (export value and edge recovery) — an implementer cannot write the marshaling struct without picking one, and RT1's "value identity on every C2 field" is unverifiable until the shape is fixed. Already tracked as A7 (Pending) — but A7 remains unresolved in this record, so it blocks coding today, not just Stage 6.

2. [High] 0021:C4 — C4 requires refusing with `graph-export-too-large` "when the traversal is incomplete under the published node ceiling (`reach.go::reach` returns `complete == false`)", but `reach` is unexported (package-private in internal/graphlint) and the only public entry point, `graphlint.Reach(m *table.Model) []Node`, discards `complete` entirely (`nodes, _ := reach(m); return nodes`). Every existing caller (lint's analysis.go:46 and ~20 test files) uses `Reach`, not `reach`. The Existing Infrastructure Audit table says to "Add a sibling exported function" for edge recovery but never names or specifies a second exported symbol (or a new `Reach` return signature) that also surfaces `complete`. First-hour question: what is the actual Go signature the `graph` verb calls to get `(nodes, edges, complete)` together — is `graphlint.Reach`'s signature changing (breaking ~20 test call sites), or is there a new `graphlint.ReachWithEdges` / similar that C4 and the Existing Infrastructure Audit never name?

3. [Medium] 0021:C1 — C1 says `--flow` alone yields `flag-invalid-value` "(this build resolves no ids)", mirroring `lint.go::runLint`'s arm set "verbatim." lint.go's own Long help text documents `--flow` as "reserved; this build resolves none — use --model," consistent with C1. However C1 does not say whether the `graph` verb's `--flow` flag description/help text must also carry that same "reserved" wording, or whether it can diverge from lint's since it's a distinct flag instance on a distinct command. Blocks Phase 2 (CLI verb and gateway wiring) — an implementer copying lint's flag definitions verbatim vs. writing new ones with equivalent semantics needs to know whether "mirrors the arm set" extends to help text, not just error codes.

4. [Medium] 0021:C5 — C5 states the verb offers "NO caller-controlled projection of its success payload" and if a later revision adds one it must conform to "JDR 0002 §D1." This RDR's own References/edges do not resolve or anchor a JDR 0002 (edges show `0021:D-identity → 0002` as RDR 0002, not a JDR); there's no `source-anchor` or resolvable reference for "JDR 0002 §D1" in the edges facet. First-hour question: is "JDR 0002" a distinct joint-decision-record from "RDR 0002" (Transition table as reviewable data) cited elsewhere in this same document? If they are different documents with colliding numbers, an implementer chasing this forward-compat clause cannot find it.

5. [Low] 0021:C2 — C2 requires the `domain` field on a tag "exactly when `guard.AssignmentCount` reports it finite." Source confirms `func AssignmentCount(d table.TagDecl) (int, bool)` exists in internal/guard/declaration.go, so the contract is implementable as written — no ambiguity found here, noted only because it was the one A3-anchored symbol worth checking and it resolved clean (not a defect).

---

§return-packet
verdict: NEEDS_DECISION
blocking: yes
evidence_paths: [/Users/cwensel/sandbox/newcoinc/intrastate/docs/rdr/0021-lint-normalized-graph-export/evidence/3amigo/persona-2-implementer.md]
changed_paths: []
next_action: Resolve A7 (fix values shape in C2) and name the exported symbol/signature C4's completeness check calls, before Phase 1 coding starts.
summary_50w: Two high-severity gaps block coding: C2's reach.values shape contradicts itself between exhibits (tracked as open A7), and C4's incompleteness check cites private reach.go::reach with no exported symbol surfacing complete. Two medium items: C1 silent on --flow help-text mirroring scope; C5 cites unresolved "JDR 0002 §D1."

Ledger:
1. 0021:C2 — reach.nodes[].values shape contradicted (object vs array) between Illustrative Code and A5 fixture; open as A7, blocks Phase 1 struct definition.
2. 0021:C4 — completeness check needs reach.go::reach's private `complete` bool but graphlint.Reach (the only exported entry point) discards it; no named exported symbol exposes it.
3. 0021:C1 — unclear whether "mirrors lint's arm set... verbatim" extends to --flow's reserved-flag help text or only to error codes.
4. 0021:C5 — forward-compat clause cites "JDR 0002 §D1," unresolved/unanchored in this record's edges; ambiguous vs. the RDR-0002 also cited throughout.
5. 0021:C2 (informational, not a defect) — guard.AssignmentCount signature matches C2's domain-inclusion test as written; verified clean.
