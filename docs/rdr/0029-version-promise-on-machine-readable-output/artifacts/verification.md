
## Phase 3b — Adversarial

Three most likely failure modes, anchored in `0029:§failure-modes`, each
with the test that catches it. All three tests fail against the current
implementation.

ADV-1 (F2, "Fails silently": an `info` code is promoted to `blocking`
without a release note; "the only way to diagnose it is to diff the
severity constants between two builds"). `internal/cli/lint.go:249`
(`nearMissFindings`) hardcodes `Severity: graphlint.SeverityInfo` at the
emit site. The code it stamps, `table.AdvisoryNearMiss`
(`reserved-tag-key/near-miss`), is in neither `graphlint.BlockingCodes()`
nor `AdvisoryCodes()` — the two seams `renderVocabularySnapshot` ranges
over — so it has no row in `internal/graphlint/testdata/vocabularies.txt`.
Editing that literal to `SeverityBlocking` changes what a consumer's CI
sees while leaving the committed snapshot byte-identical. C3 backs its
disclosure obligation with that snapshot and concedes exactly ONE
uncovered row ("the CLIError `code` row has none buildable (A5) … the one
promotion-shaped event this mechanism cannot catch"); this is a second,
unconceded one. Caught by
`internal/cli/schema_adversarial_0029_test.go::TestAdv0029_EveryWireSeverityIsCarriedByTheSnapshotRange`.

ADV-2 (F3, "Fails silently": a new machine-readable surface ships without
a tier assignment). `table.AdvisoryNearMiss` crosses the wire as
`data.findings[].code` on the SUCCESS envelope, yet belongs to no
vocabulary carrying an enumeration seam: `internal/table` owns it as a
lone exported const and exports no accessor over the advisory-rule
vocabulary. REQ-29 fixes the scope ("A vocabulary is tiered wherever it is
EMITTED, including on the fields of a `findings[]` element … an unassigned
machine-readable surface is a defect") and REQ-30 obliges the seam. C4's
census covers the graph-lint blocking/advisory code sets and the CLIError
`code` registry — this code is in none of them, and it is not among the
two namespaces C4 deliberately leaves untiered (`findings[].class`,
`data.dispositions`). Caught by
`internal/cli/schema_adversarial_0029_test.go::TestAdv0029_TheEmittedAdvisoryCodeVocabularyHasATierSeam`.

ADV-3 (F3, same bullet — "Diagnosed by auditing C4's list against the
emitted vocabularies"; this test is that audit). C4 tiers "the severity
vocabulary (`blocking`, `info`)" as `frozen` and `graphlint.Severities()`
is its seam, but the rendered snapshot has no `[graphlint.Severities]`
section. The existing `[graphlint.FindingSeverities]` section is a
different subject — the code→severity MAPPING — so adding a third member
to `severities`, a frozen set gaining a member, leaves the snapshot
byte-identical until some code is also moved onto it. C3 claims the check
"reaches every vocabulary that HAS a seam"; this one has a seam and is
unreached. Caught by
`internal/graphlint/severity_snapshot_adv_0029_test.go::TestAdv0029_TheFrozenSeverityVocabularyIsRecordedInTheSnapshot`.

## Phase 3a — CoVe

Independent pass: inputs derived from the RDR record and `req-list.md` only
(58 REQ-N plus REQ-MVV); Phase 1 tests, `coverage.md` and `deviations.md`
unread. Each REQ was given an input that would make a correct implementation
visibly violate it, and that input was run against the built tree
(`go build ./cmd/...`, branch `worktree-rdr-0029`).

Representative probes and observed behaviour:

- REQ-1/5/6/8/42: `lint --model models/examples/pricing-decision-table.toml --as=json`
  emits `{"type":"ok","schema_version":"0.1","data":{...}}`; `version --as=json`
  emits `data.version` with `schema_version` still top-level only. Zero-value
  marshal of both records keeps the key (`{"type":"ok","schema_version":""}`,
  `{"code":"x","message":"m","schema_version":""}`), so it is not `omitempty`.
- REQ-2/3/4: `clierr.SchemaVersion` is the one home; `clierr` and `respond`
  import `internal/version` nowhere.
- REQ-7/24/43/49: blocking run over `internal/table/testdata/pos-near-miss-folded.toml`
  exits 2 with `code=graph-lint-failed` and no other code, no `type` key,
  `findings` the sole structured field, exactly one `schema_version` in the
  document and none nested under a finding.
- REQ-10/44: no `DisallowUnknownFields` in `flow_input.go`; a `flow resolve`
  envelope carrying `schema_version` — and the same with an injected unknown
  top-level key, and the bare-`data` arm — all still decode through
  `flow set-state --plan` at exit 0.
- REQ-31/32/33/34/35/38: `Types=[ok]`, `Levels=[note warning]`,
  `ExitCodes=[0 1 2 3 130]`, `UnknownReasons=[absent uncomparable not-evaluated]`,
  `resolve::Blocks`, and both `Operators` mirrors identical by value and order.
- REQ-11/21: `docs/cli-output-contract.md` carries the tier table for every C4
  vocabulary with its seam and the untiered pair; `internal/graphlint/testdata/vocabularies.txt`
  records seam members and per-code severities.
- REQ-16/56: `graph-coverage-closed-by-escape`, `ClosedByEscape` and `closedBy`
  survive untouched in production code.
- REQ-46/REQ-MVV: `lint --model models/examples/release-grammar.toml --as=json`
  exits 0, `type` still `ok`, findings carry `"severity":"info"`.
- REQ-45: the 0023 golden differs from its predecessor by exactly the one
  `schema_version` key.
- REQ-51/53: the three named exact-value/cardinality assertions are retired
  from `findings_0006_test.go`; no cardinality assertion survives over any
  append-only or growing seam.

No input produced a REQ violation attributable to this pass.

## Verdict — clean
