RDR: 0021-lint-normalized-graph-export | phase: triage — close-ledger | state: COMPLETE
last: roborev triage swept the branch and closed its ledger; four IN-SCOPE findings fixed in 3a1b722, REQ-18 demote completed in 5ea3330/2acdceb; full suite green
blocker: none
changed: internal/cli/graph_dot.go, internal/cli/graph_document.go, internal/graphlint/export.go, internal/cli/graph_probe_0021_test.go, internal/cli/graph_document_0021_test.go, internal/cli/graph_dot_0021_test.go, internal/graphlint/graph_export_0021_test.go, artifacts/triage.md, artifacts/req-list.md, artifacts/coverage.md
validate: go test ./...
baseline: GREEN — `go test ./...` exit 0 at 2acdceb (parent re-ran the baseline at b409d65 before triage; the fixup re-ran it green)
next: land — `/rdr-implement-triage --resume 21 --close-and-flight` (Phase 3 §land-ship onward)
reads: artifacts/triage.md, internal/graphlint/export.go::joinedEdgesFrom, internal/cli/graph_dot.go::nodeMatchesInitial, internal/cli/graph_document.go::graphReachDocOf
session: 2026-09-13 | post-lock fixup, then triage
artifacts: req-list.md impact.md coverage.md verification.md deviations.md triage.md

Author decisions: all four resolved before triage (deviations.md carries each
on its line-leading `Status:`); `open_author_decisions = 0` at the gate.

- D2 (REQ-109, TEST-FIXTURE) — loop 1's needles became the ESCAPE SEQUENCES
  the REQ forbids, not the raw `<` `>` `&` the fixture authors.
- D3 (REQ-67, TEST-FIXTURE) — the scanning file is excluded from its own sweep.
- D4 (REQ-85, TEST-FIXTURE) — TestReq14's REQ *is* the 0/2 exit mapping, so an
  exit-only oracle is correct there; exempted by name.
- D6 (SPEC-UNDER) — closed WITHOUT an author ruling: the locked C2 already
  predicates `domain` on AUTHORED members, which is D6's own second remedy.
  D6's two residuals (the `graph_document.go:51-54` doc comment and
  TestReq20And28's unasserted arms) were fixed in `b409d65`.

## Triage outcome

Five of the six pre-triage commits carried NO roborev auto-review; they were
backfilled (jobs 7792-7796) before triage ran, which is what gave the sweep a
real window. Six findings surfaced, two dropped as superseded-at-HEAD, four
routed IN-SCOPE and all four fixed in `3a1b722`. No finding routed
OUT-OF-SCOPE, so no kata was filed and `batch:rdr-0021` has zero children.
Each fix was proven non-vacuous by reverting the production files with the new
tests kept and confirming each regression failed with the diagnosed symptom.

- **RT1 round-trip oracle** (`graph_document_0021_test.go:915`) was a
  self-consistency tautology comparing `f(x)` to `f(x)`; it now asserts the
  decoded document against the fixture's AUTHORED values.
- **Initial-node marker** (`graph_dot.go`) — `nodeMatchesInitial` skipped a
  cleared declared key, so a successor took the `initial` marker; a missing
  declared key now returns false.
- **Empty owned `set`** (`graph_document.go`) rendered `"values":{"key":null}`
  against C2/REQ-94's "never `null`"; routed through `stringsOrEmpty`.
- **Edge/successor join asymmetry** (`graphlint/export.go`) — edge recovery
  resolved successors through `indexOf` while the fixpoint joins by presence
  footprint, orphaning merged nodes (3 of 6 nodes with zero inbound edges on
  the shipped wire). New `joinedEdgesFrom` mirrors `successorsOf`. This defect
  appeared in no per-commit review; only the cross-commit range net found it.

`verification.md`'s Phase 3b entry listing the join asymmetry under
"Hypotheses probed and REFUTED" is superseded for this defect: its lemma (the
`indexOf < 0` branch is dead) still holds, but its oracle "every node retained
an inbound edge" is exactly what the cited cycle breaks.

REQ-18 was demoted to `req-list.md`'s `## EXCLUDED` section and its empty
`coverage.md` row removed, on deviation D1's existing rationale. Note that
`deviations.md`'s D1 entry still carries its pre-demote framing; it is left
as authored rather than amended.
