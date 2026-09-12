Model: claude-opus-5

# Grounding (S-00) — codebase claim sweep, cli/0029

Scope: 36 `source-anchor` edges (all `resolved:true` — the anchored half needs no
finding) plus the unanchored prose claims: counts, "no existing X", reachability,
and the inverse sibling-path check. Three delegated source sweeps; every verdict
below read on `main`.

## Findings

**GF1 — REFUTED. A4's vocabulary census is incomplete: two emitted closed
vocabularies carry no tier.**

`clierr.Finding` marshals `Operator string \`json:"operator,omitempty"\`` and
`Block string \`json:"block,omitempty"\`` (`internal/cli/clierr/clierr.go:332,334`).
Both ride `findings[]` on the `--as=json` wire, populated at
`internal/graphlint/analysis.go:152`, `internal/graphlint/coverage.go:177,336`,
`internal/graphlint/groups.go:354` and `internal/cli/flow_exec.go:788`.

Both are closed sets, and the source says so in its own words:
- `operator` — 8 tokens from `internal/guard/grammar.go::Operators`, whose doc
  comment reads "returns this RDR's **closed** operator vocabulary"
  (`eq,in,lt,lte,gt,gte,exists,contains`); `KnownOperator` adds "The set is
  closed".
- `block` — 3 tokens from `internal/resolve/guard.go::Block`
  (`all`/`unless`/`match`), the type fixed by `0007:C1`, `BlockMatch` added by
  `JDR 0001 §D12`.

Neither appears in C4. The census enumerated top-level and `data.*` fields but
not the fields INSIDE `findings[]` elements — it caught `findings[].class` (and
correctly ruled it tier-less author surface) while missing its two siblings in
the same struct. Independent count: 15 emitted string vocabularies, not 13.

This is A4's own "If wrong" case: "an unassigned surface has no stated promise —
the exact gap this RDR exists to close, reappearing inside the fix."

**GF2 — REFUTED. A3 understates the strictness risk: a golden-file test breaks
on C1's added field.**

A3 concludes "the additive field is safe for the strictest consumer that exists —
and `planEnvelope` is the one this RDR must not break," and `0029:S3` calls
`planEnvelope` "the only strictness risk the repo actually holds."

`internal/cli/flow_mvv_0023_test.go` compares a full terminal
`flow resolve --as=json` envelope against a checked-in golden by STRING EQUALITY
(`got != strings.TrimRight(string(want), "\n")`), golden at
`docs/rdr/0023-resolve-envelope-projection/artifacts/mvv-step1-default-golden.json`
(`{"type":"ok","data":{…}}` — no `schema_version`). Its own comment states the
intent: "Byte-identity is still asserted over the whole record."

A non-`omitempty` `schema_version` on the `ok` envelope — exactly what C1
mandates — fails that assertion on the commit that adds it. A3's claim is scoped
to DECODERS and holds there (no test uses `DisallowUnknownFields`); the tightest
coupling in the repo is not a decoder but a golden file, so the frame, not the
decoder finding, is what is wrong.

Not a refutation of C1: the golden is a 0023 artifact recording a pre-change
baseline, and updating it is the correct response to an intended additive change.
It is a missing IMPLEMENTATION OBLIGATION and a missing test row.

**GF3 — REFUTED (narrow). A2's release-trigger claim overstates.**

A2: "`.github/workflows/release.yml` triggers only on a `v*` tag push
(`on.push.tags: ["v*"]`)." The workflow also carries `workflow_dispatch`
(`.github/workflows/release.yml:14-16`), so "only" is false.

A2's load-bearing claim is untouched: `git tag --list` is empty and
`git describe --tags` fails (exit 128), so no release has been cut and the whole
`0.x` window is open. Only the word "only" is wrong.

**GF4 — REFUTED (narrow). A3's `decodeStrict` call-site count is wrong.**

A3 says "its **four** in-repo call sites are all TOML, reached via
`internal/table/load.go::Load`." `decodeStrict` has exactly **one** call site:
`internal/table/load.go:57`. The other `rg` hits are comment references in four
test files, not calls. The load-bearing point (input-path-only, TOML, via `Load`)
holds; the count does not.

## Confirmed (no finding owed)

All 36 `source-anchor` edges resolve `true`. Additionally read on `main` and
confirmed exact:

- `taxonomy.go::severityFor` / `IsBlocking` as described; no exhaustiveness
  switch and no cardinality constant in `internal/graphlint/` or `internal/cli/`.
- `AdvisoryCodes()` has exactly **11** call sites, all iterating or membership-
  testing; none asserts cardinality or ordinal position.
- `TestReq74_TheAdvisoryTierIsClosedAtExactlyFourMembers` is the SOLE assertion
  breaking on a fifth advisory code; its 4-element `want` literal is at
  `findings_0006_test.go:288-293`, exactly as cited.
- All five newly-assigned A4 surfaces exist as described. `cmdbind` does reject a
  stranger verdict (`internal/cli/cmdbind/cmdbind.go:1245`). Nuance: that
  rejection lives on the external gate-envelope decode path only; the internal
  `flow_exec` switches have no default-reject arm. The claim stands — cmdbind is
  the constraint that stops the verdict set appending.
- The `flow next` unknown-`reason` set (`absent`/`not-evaluated`) is genuinely
  distinct from the `graph-unprovable-coverage` `reason` set.
- `clierr::ExitCodeFor` exists. `graphlint::AggregateCode` = `graph-lint-failed`.
- `Categories()` returns exactly **40** members, consistent with 25+3+6+6.
- Both `closed`-sense comments exist verbatim (`category.go:89` "the closed
  load-category set"; `category.go:132` "The list's size is not a contract"),
  and `taxonomy.go:37-38` "The tier is CLOSED at these four (`0006:C17`)" — C2's
  retirement of the word has real, located targets.
- `Report.Blocking()`/`Advisory()` partition on a two-valued severity; only the
  blocking half reaches `respond.Fail` under `AggregateCode`.
- A2's version facts: no tag ever cut; `version::resolve` pinned to `"dev"`.
- `docs/cli-output-contract.md` has 0 occurrences of release/upgrade/breaking.
- `internal/cli/lint.go` registers exactly `--model` and `--flow`.
- `table::CatUnsupportedVersion`'s message is verbatim as quoted, including the
  odd in-production "this RDR accepts version 1 only" phrasing.

## Inverse sibling-path check

C1 adds a new envelope discriminator. Searched for an existing output-schema
version across the Go sources (`schema_version|format_version|apiVersion|
schemaVersion`): **searched, none exists**. The only version-shaped JSON field is
`internal/version/version.go:61`'s `json:"version"`, the binary's build identity —
a different surface, and the one the RDR's Naming decision already rejects reusing.
The `[model] version` input discriminator is named and dispositioned in
Load-Bearing Decisions. No sibling path already makes this decision.
