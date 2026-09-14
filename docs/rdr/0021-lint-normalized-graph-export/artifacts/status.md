RDR: 0021-lint-normalized-graph-export | phase: 3c — fixup | state: COMPLETE
last: rebased onto main @86f8c24 (locked Final record); four author decisions resolved; D2/D3/D4 test fixes + D6 residuals applied; full suite green
blocker: none
changed: internal/cli/graph_document.go, internal/cli/graph_document_0021_test.go, internal/cli/graph_stability_0021_test.go, artifacts/deviations.md, artifacts/coverage.md
validate: go test ./...
baseline: GREEN — `go test ./...` exit 0 (was red @fc1cbbb on REQ-109/D2, REQ-67/D3, REQ-85/D4)
next: run 3d on the orphan check, then the gate — `/rdr-implement-triage --resume 21 --close-and-flight`
reads: internal/cli/graph_document.go:45-60, internal/cli/graph_document_0021_test.go:134-290, internal/cli/graph_stability_0021_test.go:154-240
session: 2026-09-13 | post-lock fixup
artifacts: req-list.md impact.md coverage.md verification.md deviations.md

Author decisions, all four resolved (deviations.md carries each on its
line-leading `Status:`):

- D2 (REQ-109, TEST-FIXTURE) — loop 1 scanned for the RAW characters `<` `>`
  `&` while the fixture authors `a<b`/`x&y`, so it refused exactly what loop 2
  requires; no document could pass both. The needles are now the ESCAPE
  SEQUENCES the REQ actually forbids.
- D3 (REQ-67, TEST-FIXTURE) — the self-scan's three needles live in its own
  source, so the file matched itself. The scanning file is excluded from its
  own sweep; the other files stay covered.
- D4 (REQ-85, TEST-FIXTURE) — TestReq14's REQ *is* the 0/2 exit mapping, so an
  exit-only oracle is correct there. Exempted by name, with the reason inline.
- D6 (SPEC-UNDER) — closed WITHOUT an author ruling, and this is the one worth
  reading twice. The entry was authored against the pre-narrowing REQ-20
  ("present exactly when `guard.AssignmentCount` reports it finite") on a
  branch forked at bd7a026, eight commits before the record was refined. The
  locked C2 now predicates `domain` on AUTHORED members and states the
  carve-out explicitly — which is D6's own second remedy. The shipped gate
  (`finite && len(decl.Domain) > 0`) and the committed golden already
  implement it, so nothing in production moved.

D6's two residuals, both assigned to the implementation by A3, are fixed here:

- the `graph_document.go:51-54` doc comment no longer asserts the refuted
  "the finite arm always carries at least one member";
- TestReq20And28 (renamed `…WhenAuthored`) now asserts every arm. The fixture
  gained a `bool` and a bounded `int`, and the test pins `domain` ABSENT for
  the three finite-but-member-less kinds (`bool`, bounded `int`, member-less
  `enum`) alongside the existing enum-with-members PRESENT and `scalar` ABSENT
  arms. That is the coverage gap S2 required on both sides of the rule.
