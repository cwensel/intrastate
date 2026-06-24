# Cluster 0001-0006 Whole-Set Critique

Date: 2026-06-24

## Scope

Cluster members:

- 0001-resolution-kernel
- 0002-transition-table-as-reviewable-data
- 0003-guard-predicate-exhaustiveness
- 0004-accessor-execution-safety-model
- 0005-skill-integration-cli-contract
- 0006-graph-lint-authority-and-guarantees

Each member was Final at the start of this pass, and none was Implemented.

## Three Most Likely Inter-RDR Failure Modes

1. `flow next` overpromises legal outcomes.
   - Root cause: RDR 0005 says `flow next` returns the legal recognized-outcome
     alphabet while also saying it must not evaluate missing guard or accessor
     facts.
   - Symptom: skills present outcomes that later refuse in `flow resolve`,
     making constrained decoding feel nondeterministic.

2. The normalized model becomes a second source of truth.
   - Root cause: RDR 0002 makes sparse TOML authoritative but also requires
     normalized candidate rows, deterministic dumps, source locators, and
     parse-normalize-dump value identity.
   - Symptom: lint and resolver pass on normalized fixtures while authored TOML
     diagnostics point at weak or stale source identities.

3. Lint authority ships locally but not as a real gate.
   - Root cause: RDR 0006 depends on CI invoking production
     `intrastate lint`, and that proof remains an MVV obligation.
   - Symptom: broken graph models merge because package tests pass but
     `make check` or CI never exercises the production command.

## Most Likely Rewrite

RDR 0005 is the most likely rewrite. Its scalar-only MVP grammar refuses
duplicate `--tag` names until structured set literals exist, while RDR 0003
requires set-valued tags and `contains`. Real skill calls are likely to need
typed or bulk input earlier than a later structured input file.

## Cross-Cutting Assumption Most Likely To Fail

"Value identity is enough; byte/source identity is out of scope." The cluster
uses that stance for resolver replay, parse-normalize-dump, accessor
write-read safety, and request-level CLI determinism. It is coherent, but it is
fragile once users debug real TOML, CLI JSON, and artifact writes.

## Premortem

Six weeks after implementation, `intrastate flow next` is considered
unreliable by skill authors. It returns outcome alphabets that are technically
model-derived, but users hit `flow resolve` refusals because guards were
summarized, not evaluated. Meanwhile `intrastate lint` exists but was added as
a package-level test helper before the root command and CI path were finished,
so invalid transition fixtures slipped through review. Debugging is painful
because normalized rows have stable ids, but source spans are optional and the
CLI grammar cannot express the set-valued tags needed to reproduce a failing
`contains` predicate. The team rewrites RDR 0005 first to add structured input
and sharper `next` wording, then tightens RDR 0006 CI acceptance.

## Acceptance Tests

- `flow next` fixture where one outcome is possible only with missing guard
  facts: assert output labels it unresolved or conditional, then `flow resolve`
  without those facts refuses with the expected code.
- CLI fixture with set-valued tag input for a `contains` guard: assert scalar
  duplicate `--tag` refusal and structured input success both behave
  intentionally.
- Parse-normalize-dump round trip over reordered TOML: assert same
  candidate-row value and actionable source rule id/span in every lint finding.
- CI-shaped test: `make check` must build and invoke root `intrastate lint`
  over legal and illegal model fixtures.
- Accessor write/read-back test: mutate an owned tag while also changing a
  protected observed tag; assert read-back mismatch, not success.

