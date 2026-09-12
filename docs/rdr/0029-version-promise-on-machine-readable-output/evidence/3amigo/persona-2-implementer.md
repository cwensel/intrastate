Model: claude-opus-5

# Persona 2 — Implementer: first-hour clarification requests

Owned set read: `0029:C1`–`C4`, `0029:D-identity`/`D-wire-byte-format`/`D-naming`/`D-selection-predicate`, and the `source-anchor` edges.
**Widened** to `§implementation-plan` (Step 1/2/3), `§testing-strategy` (`0029:S1`–`S9`), `0029:MVV`, `§existing-infrastructure-audit`, and `§capability-dependencies` — because three of my first-hour questions are about what the contracts do NOT say (who bumps the version, which of two same-named symbols is tiered, what "grep for `closed`" actually matches), and a silence has no line range. Source reads were scoped to the files the anchors name.

---

## SEV-1 — blocks Phase 1 Step 1 / Step 2 on day one

### F1. `0029:C1` — the anchor `internal/cli/flow_input.go::planEnvelope` is a struct that parses INPUT, not the output envelope; C1's closing sentence misdescribes it

C1's `§illustrative-code` tail, `0029:A3`, and `0029:S3` all rest on the claim that `planEnvelope` "already does" the `code`-presence discrimination a consumer will do over the OUTPUT envelope. It is not a function (there is no `func planEnvelope` anywhere in the tree) — it is a struct at `internal/cli/flow_input.go:338` whose job is to decode a `--plan` document read from stdin/file. Two consequences an implementer hits immediately:

- The C1 claim that it "distinguishes them by `code`'s presence" is *half* right: the struct carries `Code`/`Message` for exactly that reason, but its primary discriminator, per its own doc comment, is **presence of a top-level `type` key** as raw bytes — and its doc comment explicitly says a refusal envelope on stdin is **refused**, not parsed. So it is not the "tolerant reader" pattern C1 holds it up as.
- Worse, that same doc comment states the accepted bare-`data` arm works because "`resolvePayload`'s fourteen keys are fixed and `type` is not among them." C1 adds `schema_version` to the `ok` envelope — which is the top level, not `data` — so this particular struct survives. But no contract says whether `schema_version` rides the top level ONLY or is ever projected into `data`. **Blocks:** whether Step 1's field addition can be a pure gateway change or needs a `readPlan` audit. Ask: does `schema_version` ever appear anywhere but the top level of the two terminal records?

### F2. `0029:C1` / `§step-1` — nothing says WHO increments `schema_version`, or where the literal lives

C1 fixes the form (`MAJOR.MINOR`, string, not derived from the binary) and the *rules* for when each component moves. `0029:MVV` step 4 asserts the minor moves `"0.1"` → `"0.2"` as an observable. But no contract names the mechanism: is it a `const SchemaVersion = "0.1"` in `respond`? in `clierr`? in a third package both import (they are deliberately separate packages — `clierr` is the leaf that exists to avoid an import cycle, per its package doc)? And nothing says what enforces the bump — a human in review, a test, a lint? **Blocks:** Step 1's actual edit. If the literal is duplicated in `respond.Success` and `clierr.CLIError` they can drift, which is precisely the "a key meaning two things on one wire" defect `0029:D-naming` cites as the reason this RDR exists. Ask: one exported constant in which package, and is the bump asserted or reviewed?

### F3. `0029:C2` last clause + `0029:S9` — "no occurrence of `closed`" is not greppable as written, and one occurrence is a shipped finding-code identifier

`0029:S9` says a grep assertion "is sufficient". It is not. In the two files `§step-2` names, `closed` occurs 7× in `internal/graphlint/taxonomy.go` and 3× in `internal/table/category.go`. Among them:

- `internal/graphlint/taxonomy.go:40` — `CodeCoverageClosedByEscape = "graph-coverage-closed-by-escape"`. That is a **wire-visible finding code**, a member of the very advisory vocabulary `0029:C4` tiers `growing`. Renaming it is a breaking change to an append-only set; not renaming it means a bare grep for `closed` in that file can never go to zero.
- `internal/graphlint/taxonomy.go:56,68` and `category.go:89` use `closed` in the *compatible* sense ("closed and APPEND-ONLY", "the closed load-category set"), which `0029:§key-discoveries` identifies as one of the two live senses. These are rewordings.
- `internal/table/category.go:158` uses `closed` about a *reviewer's* mental model, unrelated to tiers.
- Repo-wide, `closed` appears 99× in non-test Go. C2 says "MUST NOT be used to describe any of these tiers, in code comments or documentation" — scope is the whole repo and the docs, not the two files `§step-2` names.

**Blocks:** Step 2's scope and S9's assertion form. Ask: is S9's guard a grep with an explicit allowlist (the code identifier, the reviewer-prose site), a scoped grep over only the two named files, or a review item? And is `CodeCoverageClosedByEscape`'s literal exempt by C4's own `growing` tier?

### F4. `0029:C4` — `findings[].operator` is anchored ambiguously; there are TWO `Operators()` in the tree

C4 tiers `findings[].operator` as `frozen`, anchoring `internal/guard::Operators` and `0003:C7`. But `internal/table/model.go:96` also exports `func Operators() []string` over its own `var operators` (8 members: `eq, in, lt, lte, gt, gte, exists, contains`), mirrored from 0003 with the comment "this RDR does not mint operators and does not widen the set (`0002:C16`)". Two enumerations of the same vocabulary in two packages. C4 tiers one and is silent on the other. **Blocks:** `0029:S6` ("`frozen` sets match their declared members exactly") — the by-value assertion has two candidate subjects, and a frozen-set test that pins only one lets the other drift. Ask: does S6's frozen assertion pin both, or is `table.operators` an internal mirror out of scope? If out of scope, what keeps the mirror honest?

---

## SEV-2 — blocks a named test or a named assertion

### F5. `0029:C4` — `exit-code classes` is not a vocabulary at the granularity C4 freezes it

C4 freezes "the exit-code classes emitted by `internal/cli/clierr::ExitCodeFor`". `ExitCodeFor` (clierr.go:131) has TWO enumerations behind it: the `ErrorGroup` iota constants (`GroupSuccess, GroupWarning, GroupUserEnv, GroupEnvUnavailable, GroupInternal, GroupSignalCancel` — six, and explicitly documented "Add groups as the exit taxonomy grows") and the integer codes they map to (`0, 2, 3, 130`, plus a default `1` for a non-CLIError). `ErrorGroup` is not serialized (`json:"-"`), so it is not machine-readable output at all; the integers are. If `frozen` binds the groups, C4 contradicts the source comment inviting growth. If it binds the integers, the default `1` arm is an unnamed fifth class C4 does not list. **Blocks:** S6's by-value assertion for this row — an implementer cannot write it without knowing which list is the frozen set. Ask: freeze `{0,1,2,3,130}` as integers, or the six groups?

### F6. `0029:C1` vs `internal/cli/respond` package doc — the source says the refusal envelope IS `{"type":"failed"}`

C1 is emphatic that a refusal carries no `type` and no wrapper, "per `0005:C1`", and `§illustrative-code` bolds it. The live source disagrees in three places that are themselves machine-read documentation: `internal/cli/respond/respond.go:13` documents `{"type": "failed", ...}   # graceful failure` as the contract; line 24 says the `"ok"/"failed"` type names "are reserved"; line 49 says `Type` is "always 'ok' (Fail emits 'failed')". The *behavior* matches C1 — `Fail` calls `clierr.EmitJSON` which marshals the bare `*CLIError` — but the package doc is stale, and `internal/cli/flow_input.go:361` (readPlan's doc) ALSO says a lost pipe "carries `{"type":"failed",…}`". No implementation step names these. **Blocks:** `0029:S2` ("no `type` key appears") — a reviewer reading the package doc will read S2 as a regression, not a confirmation. Ask: does Step 1 or Step 2 also correct `respond`'s package doc and `readPlan`'s comment, and is `docs/cli-output-contract.md` similarly stale on this point?

### F7. `0029:§step-1` / `0029:S3` — the 0023 golden re-capture crosses an RDR boundary with no stated authority

Step 1 requires re-capturing `docs/rdr/0023-resolve-envelope-projection/artifacts/mvv-step1-default-golden.json`, asserted by byte-identity at `internal/cli/flow_mvv_0023_test.go:44`. Both the file and the assertion exist as described — this is accurate. What is missing is whether editing another RDR's checked-in artifact is permitted, and whether 0023's own record needs a note. `§prerequisites` lists the `0006:C17` amendment (with the exact `want` literal line range) as a prerequisite but says nothing about 0023. **Blocks:** whether Step 1 is one commit or a cross-RDR coordination. Ask: is the 0023 golden re-capture in scope as a mechanical artifact refresh, or does it need the same "ruled" treatment A1 got?

### F8. `0029:C2` `append-only` forbids cardinality assertions, but `0029:S6`'s named pattern asserts by value

C2's `append-only` clause says a consumer "MUST NOT assert on the set's cardinality". S7 correctly routes append-only sets to membership-and-uniqueness (`TestReq79_TheCategorySetIsAppendOnlyAndDuplicateFree`). But `0029:S8` keeps `TestReq74_TheAdvisoryTierIsClosedAtExactlyFourMembers` alive, merely "updated", for a vocabulary C4 moves to `growing` — which C2 defines as "as `append-only`, and additionally…", i.e. cardinality assertions are equally forbidden. Updating a test whose NAME and whose `want` literal both encode "exactly four members" to a set that may grow reproduces the defect at five. **Blocks:** the shape of S8's edit. Ask: is `TestReq74` rewritten to membership-and-uniqueness (and renamed), or deleted, or does it keep a by-value `want` that every future advisory code must edit?

---

## SEV-3 — a decision an implementer would otherwise make silently

### F9. `0029:C4` — `internal/version`'s payload takes no tier and is not listed among the deliberate no-tier surfaces

C4 closes with "an unassigned machine-readable surface is a defect", and names exactly two deliberate exemptions (`findings[].class`, `data.dispositions`) plus three free-text `Finding` fields. The `version` verb emits a payload under `--as=json` (`internal/cli/version.go`, `internal/version/version.go::resolve`, which `0029:A2` anchors). `0029:D-naming` reasons about it directly — "the `version` verb's payload already means build identity" — so it is in the author's view. It appears in neither the tier list nor the exemption list. **Blocks:** nothing hard, but under C4's own closing sentence the implementer must either tier it or declare it exempt on day one. Ask: what tier do the `version` payload's field names take?

### F10. `0029:D-identity` — "two envelopes are the same shape iff they share a `schema_version` major" is unimplementable during `0.x`

D-identity makes the major the identity key. C1 says the series starts at `"0.1"` and that "while its major is `0` the schema is explicitly unstable and MAY change incompatibly in any release". So through the entire `0.x` runway every envelope shares major `0` and D-identity declares them all "the same shape" while C1 declares they may not be. C1 supplies the resolution (the minor still moves, so a consumer "can detect movement") but D-identity does not carry the qualifier and, read alone, is the rule a consumer library would implement. **Blocks:** what a consumer-side conformance helper compares during `0.x` — major only, per D-identity, or the full string. Ask: should D-identity read "major, or the full `MAJOR.MINOR` while major is `0`"?

### F11. `0029:C1` — the tolerant-reader rule is stated as a consumer MUST with no CLI-side observable

C1: "A consumer MUST ignore object properties with unrecognized names, and MUST NOT reject an envelope reporting an unsupported major" [sic — reject]. `0029:MVV`'s `disposition` table marks both rows "n/a — consumer side". `0029:S3` tests the one in-repo decoder as a proxy. That is the honest available test, but it means C1's central compatibility clause has no assertion in this repo that survives the day `planEnvelope` is refactored. **Blocks:** whether `docs/cli-output-contract.md` (Activation Step 1) is the only artifact carrying this clause, and whether anything checks `make docs-check` covers it. Ask: does the contract doc get a machine-checkable conformance fixture, or is the consumer obligation prose-only by design?

### F12. `0029:§capability-dependencies` — `graphlint.BlockingCodes()` is cited as an existing enumerable surface but C4 tiers "the graph-lint BLOCKING finding codes" without anchoring it

The capability table names `graphlint.BlockingCodes()`; C4's `append-only` bullet names "the graph-lint BLOCKING finding codes" with no symbol anchor, while the neighbouring entries all carry one (`internal/table::Categories()`, `internal/resolve::Block`). **Blocks:** which symbol `0029:S7`'s membership-and-uniqueness assertion targets for this row. Minor, but it is one of the six `append-only` rows S7 must cover.
