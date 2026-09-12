Model: claude-sonnet-5

## Queries run

1. `arc corpus list` — tier 1 — hit (corpora enumerated: Claude, DevRef, others).
2. `arc search semantic --corpus DevRef --limit 8 --json "API versioning policy breaking change additive fields JSON output stability"` — tier 1 — miss for CLI-specific policy (returned generic REST/microservices API-versioning book excerpts, nothing about lint codes, refusal codes, or CLI JSON envelopes). Not re-run at this tier per budget rule.
3. `ls /Users/cwensel/sandbox/langref` — tier 2 setup — hit: confirmed presence of `eslint`, `golangci-lint`, `gh-cli`, `opentofu` (terraform-format heir), `squawk`, `sqlfluff`; absence of plain `terraform`, `clippy`/`rustc`, `shellcheck`, `ruff`, `staticcheck`, `npm`, `kubectl`, `cargo` checkouts.
4. `mcp__semble__search "semver policy for new lint rules breaking vs non-breaking"` on `/Users/cwensel/sandbox/langref/eslint` — tier 2 — miss (returned rule implementation code, not policy docs; semble here indexes code, not prose policy).
5. `mcp__semble__search "new default enabled linter breaking change policy release notes"` on `/Users/cwensel/sandbox/langref/golangci-lint` — tier 2 — miss (returned implementation of linter-diffing code, not the policy statement).
6. `grep` sweep of `/Users/cwensel/sandbox/langref/golangci-lint/CHANGELOG.md` for "breaking|new linter|default.*enable|major version" — tier 2 — hit: confirmed "New linters" appear as regular numbered items inside non-major releases (e.g. v2.12.0), and "⚠️ Breaking change" is used as an explicit inline marker for specific linter-internal rule reorganizations, not for linter additions.
7. `find`/`grep` of `/Users/cwensel/sandbox/langref/opentofu/website/docs/internals/json-format.mdx` — tier 2 — hit: full verbatim JSON format-version policy (see A-3 below). OpenTofu retains Terraform's original doc language.
8. `find`/`grep` of `/Users/cwensel/sandbox/langref/eslint/docs/src/maintain/manage-releases.md` — tier 2 — partial hit: confirms `semver-minor` / `semver-major` terminology is load-bearing in eslint's own process, but the full patch/minor/major bullet classification for rules lives on a different, not-locally-checked-out page.
9. WebFetch `https://eslint.org/docs/latest/contribute/report-a-bug#semantic-versioning-policy` — tier 3 — miss (404; URL structure changed).
10. WebFetch `https://golangci-lint.run/product/roadmap/` — tier 3 — miss (wrong path).
11. WebSearch "ESLint semantic versioning policy semver-major semver-minor new rule" — tier 3 — hit (secondary-source paraphrase with fragments; used to locate real page, and independently corroborated by typescript-eslint.io in query 15).
12. WebSearch "golangci-lint breaking change new default linter major version policy" — tier 3 — hit, led to correct URL `https://golangci-lint.run/docs/product/roadmap/`.
13. WebSearch "rustc clippy lint levels allow warn deny future-incompat policy" — tier 3 — hit, led to RFC 2476.
14. WebFetch `https://golangci-lint.run/docs/product/roadmap/` (corrected path) — tier 3 — hit: verbatim quote obtained (see B-1).
15. WebFetch `https://rust-lang.github.io/rfcs/2476-clippy-uno.html` — tier 3 — hit: verbatim quote on lint-level categories (see B-3), but RFC does not itself state a future-incompat mechanism for Clippy (that's a rustc-specific mechanism; see Gaps).
16. WebFetch `https://eslint.org/docs/latest/maintain/manage-releases` — tier 3 — miss (page exists but does not contain the patch/minor/major rule-classification policy; that content lives elsewhere).
17. WebSearch "eslint.org semantic versioning policy site documentation page rules" — tier 3 — hit (paraphrase with quoted fragments; canonical page URL not resolved before budget cutoff — see Gaps).
18. WebFetch `https://raw.githubusercontent.com/eslint/eslint.org/main/src/content/docs/latest/contribute/report-a-bug.md` — tier 3 — miss (404, wrong repo path).
19. WebFetch `https://typescript-eslint.io/users/versioning/` — tier 3 — hit: verbatim quote from an independent peer confirming "adding a rule is not breaking" / "recommended config changes are breaking" pattern (see B-2).
20. WebSearch "ruff versioning policy new rule stable preview unstable" — tier 3 — hit, led to `https://docs.astral.sh/ruff/versioning/`.
21. WebFetch `https://docs.astral.sh/ruff/versioning/` — tier 3 — hit: verbatim quotes (see B-4).
22. WebSearch "staticcheck new check policy default enabled breaking" — tier 3 — hit but no formal published policy found (see Gaps/Rejected).
23. WebSearch "shellcheck new check policy versioning breaking" — tier 3 — hit but confirms **no formal policy exists** (negative result, see B-5).
24. WebSearch "kubectl -o json output stability policy breaking change kubernetes" — tier 3 — hit but no explicit JSON-breaking-change policy found; only general kubectl conventions and unrelated deprecation-policy page (see Gaps).
25. WebSearch "gh cli --json flag output stability guarantee documentation" — tier 3 — hit but confirms **no documented stability guarantee** for `gh`'s `--json` output (negative result, see B-6).
26. WebSearch "cargo --message-format=json stability policy unstable" — tier 3 — hit, surfaced `cargo metadata --format-version` as the quotable analog.
27. WebFetch `https://doc.rust-lang.org/cargo/reference/external-tools.html` — tier 3 — hit: verbatim quote on `cargo metadata --format-version` (see A-4); explicitly found NO equivalent stated policy for `--message-format=json` itself (see Gaps).

## (A) Class findings

**A-1. Terraform/OpenTofu — JSON output format has its own independent version, decoupled from product version.**
Source: `/Users/cwensel/sandbox/langref/opentofu/website/docs/internals/json-format.mdx` (local checkout; language inherited verbatim from Terraform's own docs).
> "The output includes a `format_version` key, which has value `"1.0"`. The semantics of this version are:
> - We will increment the minor version, e.g. `"1.1"`, for backward-compatible changes or additions. Ignore any object properties with unrecognized names to remain forward-compatible with future minor versions.
> - We will increment the major version, e.g. `"2.0"`, for changes that are not backward-compatible. Reject any input which reports an unsupported major version."
Forces for intrastate: the `--as=json` envelope should carry its own `format_version` (or equivalent) independent of `intrastate`'s product version, with an explicit "ignore unrecognized keys" forward-compat instruction to consumers, and additive-field changes bump only the format's minor component — never the CLI's.

**A-2. golangci-lint — "might break your lint build" is treated as inherent to minor releases, not deferred to major.**
Source: `https://golangci-lint.run/docs/product/roadmap/` (fetched).
> "Patch release (intended to not break your lint build)" versus "Minor release (might break your lint build because of newly found issues)" versus "Major release (likely to break your lint build)."
> "We consider the removal of a linter as non-breaking changes for golangci-lint itself. No major version will be created when a linter is removed."
Forces for intrastate: precedent exists for explicitly telling consumers that *minor* releases (not just major) can surface new findings/refusals — the tool's job is to state this plainly rather than pretend additive diagnostics are "non-breaking" in the SemVer sense. This directly answers the RDR's fork: golangci-lint's policy admits new-finding-induced CI breakage is a known, named, *minor*-release-level risk category, distinct from both patch (safe) and major (removal/rename).

**A-3. ESLint — "semver-minor"/"semver-major" terminology gates what can merge during a release freeze.**
Source: `/Users/cwensel/sandbox/langref/eslint/docs/src/maintain/manage-releases.md` (local checkout), line 49.
> "After the patch release has been published (or no patch release is necessary), close the release issue and inform the team that they can start merging in semver-minor changes again."
Forces for intrastate: confirms the vocabulary `semver-minor`/`semver-major` is operationally load-bearing (gates PR merges, not just a doc label) — worth adopting the same explicit label vocabulary for intrastate's own contribution/release process, tagging "new lint code" PRs as semver-minor-equivalent by construction.

**A-4. Cargo — `cargo metadata --format-version` is a second, independently-versioned machine-readable schema (distinct from `--message-format=json`).**
Source: `https://doc.rust-lang.org/cargo/reference/external-tools.html` (fetched).
> "The format is stable and versioned. When calling `cargo metadata`, you should pass `--format-version` flag explicitly to avoid forward incompatibility hazard."
Forces for intrastate: a second real-world instance (alongside Terraform) of a CLI decoupling its machine-readable schema version from the product's SemVer, and explicitly requiring the consumer to pin/pass a version flag — reinforces that intrastate's `--as=json` envelope should require or default an explicit schema-version field rather than relying on product version.

**A-5. Ruff — new rules ship gated behind an explicit "preview" mode for at least one minor release before promotion to stable.**
Source: `https://docs.astral.sh/ruff/versioning/` (fetched).
> "Ruff uses a custom versioning scheme that uses the minor version number for breaking changes and the patch version number for bug fixes."
> "New rules should always be added in preview mode. New rules will remain in preview mode for at least one minor release before being promoted to stable."
Forces for intrastate: strongest class-level precedent for an explicit "experimental/unstable" marker gating new-finding introduction — directly answers "how new lint rules are introduced without a major bump": they are not introduced as stable at all until proven, and the CLI keeps a documented "does not affect default behavior yet" bucket for them.

## (B) Instance findings — the decisive question

**B-1. golangci-lint: new linters are minor-release items; NOT explicitly labeled breaking; enabling one by default is the actual risk point.**
Source: `https://golangci-lint.run/docs/product/roadmap/` (fetched) + local `/Users/cwensel/sandbox/langref/golangci-lint/CHANGELOG.md` v2.12.0 entry.
> "A new linter is added" [triggers a] "Minor release (might break your lint build because of newly found issues)."
CHANGELOG corroboration (verbatim from local checkout, v2.12.0, released 2026-05-01):
> "1. New linters
>    * Add `clickhouselint` linter https://github.com/ClickHouse/clickhouse-go-linter"
— filed as a plain numbered item in an ordinary (non-major) release, alongside bug fixes and option additions, with no "breaking" marker. Contrast with the same file's explicit `⚠️ Breaking change` marker used elsewhere only for *behavior reassignment* within an existing linter (revive moving checks between rules). **Conclusion: golangci-lint's own policy explicitly classifies "new linter/rule causes previously-clean input to newly fail" as a first-class minor-release risk, distinct from (and lesser than) a major-release breaking change** — this is the closest real-world precedent to intrastate's exact fork, and it comes down on the side of "additive diagnostic ≠ SemVer-major, but IS explicitly flagged as build-breaking-capable in the minor-release definition itself."

**B-2. ESLint / typescript-eslint: adding a rule is minor; changing what "recommended" flags is major.**
Source: WebSearch fragment quoting eslint.org (page URL not independently re-confirmed by direct fetch — see Gaps) + `https://typescript-eslint.io/users/versioning/` (fetched, independent peer, corroborating).
ESLint (via search-engine-quoted fragment, unconfirmed by direct fetch):
> "Minor Release... A new rule is created" / "Major Release... `eslint:recommended` is updated and may result in new errors (e.g., rule additions, most rule option updates)."
typescript-eslint (fetched, confirmed verbatim):
> "A change to the plugins shall be considered breaking if it will require the user to change their config."
— and per the fetch-agent's synthesis of that same page: adding a rule itself is not breaking, but updating what a "recommended"/preset config enables IS breaking. **Conclusion: this is the cleanest resolution of intrastate's fork found in the wild — the operator that breaks CI is not "new code exists" but "a previously-opted-in bucket (recommended/default-on) newly fires it."** The breaking edge is drawn at the *default-enablement* boundary, not at the *code's existence* boundary.

**B-3. clippy/rustc: lint-level defaults (allow/warn/deny) are the actual breaking-change gate, not the lint's existence.**
Source: `https://rust-lang.github.io/rfcs/2476-clippy-uno.html` (fetched).
> Deny (Correctness): "Probable bugs, e.g. calling `.clone()` on `&&T`." Warn (Style, Complexity, Perf): fixes that don't cause semantic changes. Allow (Pedantic, Nursery, Cargo, Restriction, Internal, Deprecated): "less consensus or unpolished lints."
> "Clippy will have the same idea of lint stability as rustc; essentially we do not guarantee stability under `#[deny(lintname)]`."
**Conclusion:** clippy explicitly disclaims stability for anyone who opts into `deny`-level enforcement — the tool tells consumers up front that turning warnings into hard failures forfeits a stability guarantee. This is a *published disclaimer* pattern intrastate could adopt for its refusal-code / lint-finding vocabulary: state plainly that `--fail-on-warn`-equivalent consumers accept new-finding risk on minor bumps, while default/non-strict consumers do not. Rustc's dedicated "future-incompat" lint mechanism (which specifically warns about upcoming *deny-by-default* promotions before they land) was referenced in search results but not independently opened/quoted at the RFC page — see Gaps.

**B-4. Ruff: new rules launch inside "preview" and cannot cause default-path failures until explicitly promoted, which is itself a minor bump.**
Source: `https://docs.astral.sh/ruff/versioning/` (fetched).
> "New rules should always be added in preview mode. New rules will remain in preview mode for at least one minor release before being promoted to stable."
**Conclusion: this is the most direct architectural answer to intrastate's fork** — Ruff resolves "does a new rule count as breaking" by making it categorically true that a stable-mode consumer's clean input NEVER newly fails from a new rule, because the new rule literally cannot fire without the consumer opting into `--preview`. Promotion to stable is then a normal, *disclosed* minor-version event (not a silent one), so a consumer that reads release notes knows exactly which minor bump might change their result.

**B-5. shellcheck: no formal published policy on new-check introduction (negative result).**
Source: WebSearch "shellcheck new check policy versioning breaking" — no CONTRIBUTING/docs page surfaced any explicit policy; the only relevant community guidance found was informal advice ("pin a specific ShellCheck version to avoid surprise build breaks when a new version with new warnings is published") and a `--list-optional` mechanism distinguishing default vs. optional checks, but no stated SemVer-breaking classification. **Absence of policy is itself evidence intrastate must not assume "everyone already does this" — some prominent peers simply don't.**

**B-6. gh CLI: `--json` has no documented output-stability guarantee at all (negative result).**
Source: WebSearch "gh cli --json flag output stability guarantee documentation" — the one explicit stability statement found is the opposite of a guarantee, for a non-json-flagged command: "the current output of `gh auth status` is formatted for human viewing and isn't guaranteed to remain stable." No equivalent statement of guarantee OR caveat was found specifically for `--json` output across `gh` subcommands. **Treated as a gap/negative result, not confirmed either way** — see Gaps.

## Rejected branches

- **kubectl deprecation policy page** (`kubernetes.io/docs/reference/using-api/deprecation-policy/`) — opened via search snippet only; explicitly scoped to "significant, user-visible behaviors" and API deprecation, not `-o json` output-shape stability specifically. Set aside as off-target for the CLI-output-format question (it's an API-deprecation policy, not an output-format-versioning policy); would need a dedicated fetch of `kubectl conventions` page to get closer, not done here (see Gaps).
- **staticcheck** — search surfaced only mechanics (which checks are on/off by default, `--list-optional`-style config) with no explicit SemVer/breaking-change statement about *introducing* new checks. Set aside as a probable non-answer rather than pursued further, consistent with budget (≤3 queries/claim already spent).
- **squawk, sqlfluff** (local checkouts present in `/Users/cwensel/sandbox/langref`) — not opened. Neither was named in the assignment's peer list; skipped to preserve budget for the explicitly named tools.
- **DevRef arc corpus** — one query run, returned only generic REST/microservices API-versioning book content with no CLI/lint-specific material. Not re-queried at this tier (budget rule: escalate on failure, don't repeat).
- **rustc's dedicated "future-incompat" lint mechanism** — referenced by WebSearch snippets but the RFC page fetched (2476-clippy-uno) is Clippy-specific and does not itself define rustc's future-incompat-report mechanism (that lives on a separate rustc dev-guide / `unstable-book` page). Not pursued further under budget; treated as a gap, not asserted.

## Gaps

- **ESLint's canonical patch/minor/major policy page** — never independently opened via direct WebFetch (two attempted URLs 404'd: `.../contribute/report-a-bug#semantic-versioning-policy` and the raw GitHub markdown path). What's reported in B-2 for ESLint itself is a **search-engine paraphrase with embedded quote fragments**, not a verbatim page fetch — treat the exact ESLint wording as unconfirmed; only the typescript-eslint.io corroboration (B-2) was independently fetched and verified verbatim.
- **kubectl `-o json` / `-o jsonpath` explicit stability policy** — no page was found or opened stating whether adding new JSON fields, or new resource print-columns, is breaking. The `kubectl conventions` page was returned in search results but never fetched. Treat kubectl's policy as **not found**, not as "kubectl has no policy."
- **gh CLI `--json` stability guarantee** — same status: not found either confirming or denying a guarantee specific to `--json` (as opposed to human-readable default output, which is explicitly disclaimed). Treat as unresolved.
- **cargo `--message-format=json` policy for the message stream itself** (as opposed to `cargo metadata --format-version`) — WebFetch of the official external-tools page explicitly returned "no explicit stability policy statement" for this specific flag. A secondary WebSearch paraphrase claimed "new fields may be added... a new version like `--message-format=json-v2` would probably be introduced," but this could not be verified verbatim on the official page and should be treated as folklore/forum content (it traces to a Rust users-forum thread, not docs.rust-lang.org) unless re-verified.
- **staticcheck and shellcheck formal breaking-change policies** — both treated as confirmed-absent (B-5) but this rests on WebSearch snippets only, not an exhaustive read of each project's CONTRIBUTING/docs directories. A true negative would require opening each project's docs tree directly; not done here under budget.
- **rustc's "future-incompat" report mechanism** — named in the assignment as a specific thing to check; not independently confirmed via a direct fetch of any rustc/cargo doc page in this pass (see Rejected branches). This is the single most notable gap given the assignment explicitly called it out by name.
