# RDR 0015 — Stage 2 grounding micro-sweep (step 7.5)

Date: 2026-08-28. One fresh-context sub-agent, anchors-only brief (no
justifying prose passed). Verdict: **clean — 16/16 CONFIRMED**, none
REFUTED, none NOT-FOUND.

## Return packet (verbatim)

```
verdict: PASS
blocking: no
evidence_paths: [internal/graphlint/analysis.go:314,
internal/graphlint/analysis.go:342, internal/graphlint/analysis.go:371,
internal/graphlint/analysis.go:391, internal/graphlint/analysis.go:411,
internal/graphlint/analysis.go:446, internal/graphlint/reach.go:90,
internal/graphlint/reach.go:449,
internal/graphlint/adversarial_0006_test.go:411,
internal/graphlint/adversarial_0006_test.go:501,
internal/graphlint/adversarial_0006_test.go:515, models/rdr.toml:127,
docs/rdr/0006-graph-lint-authority-and-guarantees/artifacts/req-list.md:123,
docs/rdr/0006-graph-lint-authority-and-guarantees/artifacts/req-list.md:285,
docs/rdr/0006-graph-lint-authority-and-guarantees/artifacts/req-list.md:301,
docs/rdr/0006-graph-lint-authority-and-guarantees/artifacts/req-list.md:319,
docs/rdr/0006-graph-lint-authority-and-guarantees/artifacts/deviations.md:317,
docs/rdr/0006-graph-lint-authority-and-guarantees/artifacts/deviations.md:352,
docs/rdr/0006-graph-lint-authority-and-guarantees/artifacts/deviations.md:370,
docs/jdr/0001-resolve-kernel-seam.md:941, docs/model-authoring.md:408,
docs/rdr/0006-graph-lint-authority-and-guarantees/artifacts/verification.md:192]
changed_paths: []
next_action: none
summary_50w: 1: CONFIRMED; 2: CONFIRMED; 3: CONFIRMED; 4: CONFIRMED;
5: CONFIRMED; 6: CONFIRMED; 7: CONFIRMED; 8: CONFIRMED (seven dropped
rows, each writes status="abandoned"); 9: CONFIRMED; 10: CONFIRMED;
11: CONFIRMED; 12: CONFIRMED; 13: CONFIRMED; 14: CONFIRMED (reach.go);
15: CONFIRMED; 16: CONFIRMED.
```

## Anchor list checked (brief items)

1. `analysis.go::checkDeadEnd` split-and-bound doctrine — CONFIRMED
2. `analysis.go::hasOutgoingOrdinaryRow` existential + residual doc — CONFIRMED
3. `analysis.go::splitNode` — CONFIRMED
4. `analysis.go::terminalKeys` — CONFIRMED
5. `analysis.go::satisfiesSomeTerminal` universal — CONFIRMED
6. `reach.go::reach` merged fixpoint writer — CONFIRMED
7. `TestAdvDeadEndExistentialOnMergedNode` pin + positive arm — CONFIRMED
8. `models/rdr.toml` seven `stage = "dropped"` rows, each also writing
   `status = "abandoned"` — CONFIRMED
9. REQ-37 / REQ-111 / REQ-117 / REQ-122 quotes in req-list.md — CONFIRMED
10. D12 status, nine-finding measurement, trigger language — CONFIRMED
11. JDR 0001 §JD-23 charter and no-unilateral-change clause — CONFIRMED
12. cli/0022:C1 population carve, cli/0022:C2 shared doctrine (via
    projector) — CONFIRMED
13. `docs/model-authoring.md` dead-end guarantee prose — CONFIRMED
14. `matchSatisfiable` existential-per-held-value (lives in `reach.go`)
    — CONFIRMED
15. 0006:C17 advisory tier closed at four; external-checker rejection
    (ALT4) — CONFIRMED
16. SC-23 census in verification.md — CONFIRMED

Note folded from item 14: `matchSatisfiable` is defined in `reach.go`,
not `analysis.go`; the record cites it without a file attribution, so
no correction was needed.
