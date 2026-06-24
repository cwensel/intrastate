# Critique

## 1. Three Likely Implementation Failures

### C1 Command authority drifts away from the resolver CLI namespace

Root cause in the RDR: the draft declares the authoritative surface as
`intrastate lint`, while RDR 0005 declares the user-facing resolver group as
`flow` and leaves graph lint out of scope. The current text never says whether
lint is a root command, a `flow lint` subcommand, or a future compatibility
alias.

Passage that enabled it: `The authoritative command surface is intrastate lint
over model files using the same graph-lint package.` The same phrasing appears
in the Normative Contracts and Validation scenarios.

Symptom the user sees: one implementer wires `intrastate flow lint` to stay near
`flow next` and `flow resolve`, while tests and docs expect `intrastate lint`.
CI fails to find the command or, worse, both command paths appear with subtly
different flags.

### C2 Machine-readable findings are promised but not made implementable

Root cause in the RDR: the draft says JSON mode carries findings through an
append-only optional `CLIError` field, but the current `CLIError` has only
`Code`, `Message`, `Param`, `Detail`, and `Hint`. It does not name the field,
the Go type, or whether graph-lint owns extending `clierr`.

Passage that enabled it: `The wire-visible error may add an optional findings
field to CLIError under the existing append-only omitempty envelope rule; each
finding record carries its own stable code and identity.`

Symptom the user sees: text mode shows a useful multiline summary, JSON mode
only has `detail`, and downstream checks cannot assert individual
`graph-overlap` or `graph-owned-before-write` records. A later refactor adds a
verb-local JSON wrapper instead of extending the shared failure envelope.

### C3 CI authority is named but not enforceable from the repo's actual gate

Root cause in the RDR: the draft repeatedly says CI is the authority, but the
implementation and validation plan do not say which repository gate changes or
what command CI must run after the binary exists.

Passage that enabled it: `CI must run that production command shape for
transition-model changes` and Phase 3's `add the CI-shaped invocation to the
validation path.`

Symptom the user sees: `make check` passes, GitHub CI passes, and an illegal
transition fixture lands because the lint command was only tested in a unit
test. Maintainers believe CI is authoritative because the RDR says so, but
`.github/workflows/ci.yml` never invokes the production command over the
transition model corpus.

## 2. Section Rewritten Within Six Weeks

`Phase 3: CLI And CI Surface` will be rewritten first. It currently compresses
three contracts into one sentence: command placement, invocation flags, and CI
wiring. As soon as RDR 0005's `flow` command group lands, this section will
need to answer whether lint lives at root, under `flow`, or both. As soon as the
first transition fixture lands, it will also need to name the actual `make`
target or workflow step that makes CI authoritative.

## 3. Assumption That Will Fail First Contact

A5 will not survive first contact with implementation: `CI can run
intrastate lint as the blocking graph-acceptance authority for transition model
changes.` The RDR verifies this with its own MVV test plan, not with an existing
repository gate or a concrete future gate contract. A production-command test is
not the same as CI authority.

## 4. Premortem

The graph-lint RDR shipped and the team believed illegal graphs were blocked
before resolver use. The implementation added a graph-lint package and a command
under `flow lint` because `flow next` and `flow resolve` were already the
resolver-facing surfaces from RDR 0005. Tests written from RDR 0006 expected
`intrastate lint`, so the implementer added a root alias late in the work. The
alias parsed a slightly different model flag, and fixture tests only exercised
the root command.

The first failure appeared during a model update that introduced two candidate
rows for the same state/outcome. `flow resolve` refused at runtime with a
multi-match error, but CI had accepted the model. The workflow still only ran
`make check`, and `make check` only ran unit tests. No production `intrastate
lint --flow rdr --model ...` invocation touched the repository's model corpus.
Maintainers searched the RDR and found the sentence saying CI was the authority,
but no exact workflow step or make target existed to enforce it.

The second failure was worse for users. JSON callers wanted to surface every
blocking invariant in one report. The lint engine collected findings with codes
like `graph-overlap` and `graph-coverage-gap`, but `clierr.CLIError` had no
typed findings field. The command squeezed the list into `Detail`, which looked
fine in text mode and was nearly unusable in JSON mode. A downstream tool could
only assert `graph-lint-failed`, not the individual invariant codes. The
machine-readable guarantee in RDR 0006 was true in prose and false on the wire.

The user journey failed exactly where the RDR promised safety. A maintainer
edited the RDR transition model, ran the documented checks, and merged. Later,
an agent asked `flow next` for legal outcomes and chose an outcome that led to
ambiguous resolver state. Runtime refusal protected the artifact from a bad
write, but the design-time guarantee was gone: the illegal graph was discovered
one call at a time.

## 5. Acceptance Tests That Would Have Caught The Failures

1. Given the resolver CLI contract defines `flow next` and `flow resolve`, when
   graph lint is implemented, then exactly one canonical command path is
   documented, tested through Cobra, and accepted in CI; any alias must call the
   same command object or same request builder.
2. Given an illegal fixture with two blocking findings, when
   `intrastate lint --as=json` runs, then the failed terminal record includes a
   typed `findings` array with both individual invariant codes, model identity,
   severity, and source identity fields.
3. Given `.github/workflows/ci.yml` and `Makefile`, when graph lint is
   implemented, then the repository gate invokes the built binary through the
   production lint command over the checked-in transition model or fixture
   corpus, not just unit tests of the engine package.
4. Given both text and JSON modes, when a fixture has multiple blocking
   findings, then text may summarize but JSON must preserve the same finding
   identities in deterministic order.
