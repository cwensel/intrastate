# Critique Premortem

## Three Most Likely Implementation Failures

1. **Unsupported model versions are accepted because `version` is present but semantically inert.**
   - Root cause in the RDR: the schema lists flow metadata with table id, version, and description, and the spike fixture includes `[model].version`, but no normative contract says what the loader does with an unsupported version.
   - Passage that enabled it: `Technical Design` item 1 names `version`; the source-layout contract names `[model]`; `Performance Expectations` and validation promote the spike fixture, whose version is `1`.
   - User symptom: a contributor changes `version = 2` while experimenting with a new field, the parser accepts the file as if it were v1, and lint/resolution fail later with misleading unknown-field or missing-tag errors.

2. **The expanded table dump is deterministic by accident, not by contract.**
   - Root cause in the RDR: the document requires deterministic candidate rows and cites a byte-identical spike hash, but the normative section does not define the dump value or the ordering dimensions that make the hash meaningful.
   - Passage that enabled it: `Normative Contracts` says rows are deterministic for table dumps; `Round-Trip / Inverse Invariants` says `parse -> normalize -> dump = expanded-table value identity`; `Performance Expectations` cites a SHA-256 over spike output.
   - User symptom: one implementation sorts rows by rule id while another sorts by expanded source locator; both can claim "deterministic," but golden tests and reviewer diffs churn when map iteration, predicate grouping, or write ordering changes.

3. **Stable CLI parse/lint failures drift into generic errors.**
   - Root cause in the RDR: it says malformed rules and validation failures must be stable CLI errors, but it leaves table-specific code names to later integration and does not require validation categories to survive as structured diagnostics.
   - Passage that enabled it: `Existing Infrastructure Audit` says table-specific codes are not present yet; `Testing Strategy` expects failures to become stable `CLIError`s only "when surfaced by CLI commands."
   - User symptom: the first implementation returns `config-invalid` or an internal error for unknown tags, unknown contexts, unsupported versions, and overlapping rows, so scripts cannot distinguish authoring mistakes from tool bugs.

## Section Rewritten Within Six Weeks

`Performance Expectations` will be rewritten first. It currently carries a byte-identical hash from the spike while the locked contracts do not define the dump's value shape, row ordering, predicate ordering, write ordering, or version behavior. The first golden-file test will force this section either to become a real dump contract or to remove the hash as non-normative spike trivia.

## Assumption Least Likely To Survive A Real User

A2 is the fragile assumption: "Row order is not part of successful edge selection." A real table author will eventually write two broad sparse rules that overlap and expect the later, more-specific row to override the earlier one. The RDR correctly rejects first-match semantics, but it does not yet make the author-facing ambiguity diagnostic concrete enough to prevent "just add priority" pressure after the first confusing lint failure.

## Premortem

The transition-model implementation shipped and failed during the first attempt to encode the full RDR prelock graph. The author copied the spike's `[model]` block, bumped `version = 2` to mark a new local shape, added a context inheritance case, and expected the dump to show the change. The parser accepted the model because `version` was decoded as ordinary metadata rather than a compatibility gate. The normalizer then treated v2-only fields as unknown or absent, so the first visible error came back as a generic validation failure instead of `unsupported transition model version`.

The second failure arrived in review. `transition dump` produced rows in a stable order on one machine and a different stable order on another because implementation sorted rule ids but not predicate keys inside inherited contexts. The developer pointed at the RDR's deterministic-language requirement; QA pointed at the spike hash; neither could say whether `unless` predicates should appear before local matches, after local matches, or in authored order. Golden tests started asserting the current string form, so later refactors changed review output even when candidate-row semantics were identical.

The third failure reached the user journey RDR 0002 was meant to improve. A flow author added a broad `continue-prelock` rule and a narrower `continue-large-prelock` rule. Both matched the same tag-set. Because the RDR correctly forbids row priority, the resolver refused, but the lint message lacked a stable category and did not name both rule ids plus their expanded predicate overlap. The author asked for first-match priority because that was easier to understand than the tool's refusal.

The named functions that failed were predictable: the future parser's `Load` accepted an unsupported `[model].version`; the future normalizer's `Normalize` produced candidate rows without a fully ordered value representation; the future dump function rendered a deterministic-looking but under-specified view; the future CLI command wrapped all of it through `respond.Fail` without table-specific diagnostic categories. The user journey that failed was "edit sparse TOML, run parse/lint/dump, understand exactly which authored rule to fix." Instead, version compatibility, deterministic dumps, and ambiguity diagnostics were all discovered by implementation taste.

## Acceptance Tests That Would Have Caught This

1. Given a transition model with `[model].version = 2`, when the loader parses it, then it refuses before normalization with a stable unsupported-version validation category and no candidate rows are produced.

2. Given two semantically identical fixtures whose context predicates, rule predicates, guards, and writes are authored in different TOML key orders, when they are normalized and dumped, then the expanded-table value is identical.

3. Given a fixture with overlapping broad and narrow rules, when lint runs, then the failure names both `(model id, rule id)` identities and reports an ambiguous-overlap validation category rather than relying on row order.

4. Given an unknown tag, unknown context, write to a non-owned tag, unknown accessor, unsupported version, and overlapping rows, when the CLI surfaces each validation failure, then each response uses `respond.Fail` with a `CLIError` whose code is stable and whose detail includes the authored source locator.
