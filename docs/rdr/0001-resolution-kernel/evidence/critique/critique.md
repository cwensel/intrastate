# Critique - RDR 0001 Resolution Kernel Contract

## 1. Three likely implementation failures

### Failure 1: refusal classes drift between kernel and CLI

Root cause: RDR 0001 requires "one typed refusal" but does not enumerate the
kernel refusal taxonomy or the stable discriminator the CLI must map.

Passages that enable it:

- `docs/rdr/0001-resolution-kernel.md:150` through `159` says illegal,
  ambiguous, incomplete, or unmodeled input returns an explicit refusal class.
- `docs/rdr/0001-resolution-kernel.md:180` through `189` says resolve returns
  exactly one typed refusal and lists refusal conditions.
- `docs/rdr/0001-resolution-kernel.md:272` says resolver-specific codes are
  added outside the kernel package.

Symptom: the kernel returns Go errors or ad hoc strings such as
`ErrNoMatch`/`ambiguous`, while RDR 0005 invents CLI codes such as
`zero-match` and `multi-match`. A skill that branches on CLI JSON cannot map the
kernel result back to the original transition failure without brittle string
matching.

### Failure 2: implementation treats "typed refusal" as an error path

Root cause: the RDR separates CLI output from the kernel, but it never states
that kernel refusals are value-level dispositions rather than Go `error`
returns. Existing `internal/cli/clierr::CLIError` and
`internal/cli/respond::Fail` are CLI-facing symbols, not kernel result types.

Passages that enable it:

- `docs/rdr/0001-resolution-kernel.md:171` through `176` says the kernel exposes
  structured success/refusal values but leaves the shape implicit.
- `docs/rdr/0001-resolution-kernel.md:408` through `411` asks Phase 1 to define
  the package boundary, result taxonomy, and entry point, but does not constrain
  what counts as the taxonomy.

Symptom: resolver tests assert `require.Error` for no-match and many-match.
Later, CLI code wraps those as `CLIError`, and a refusal becomes
indistinguishable from parser bugs, IO errors, or programmer mistakes in logs
and tests.

### Failure 3: replay validation misses the highest-risk refusal classes

Root cause: the MVV names no-match, ambiguous matches, and unmodeled outcome,
but the normative contract also names required owned state unavailable and
guard cannot be evaluated. Those two are the cases most likely to involve
accessor/guard boundary mistakes.

Passages that enable it:

- `docs/rdr/0001-resolution-kernel.md:187` through `189` requires refusal for
  unavailable owned state and unevaluable guards.
- `docs/rdr/0001-resolution-kernel.md:402` through `406` only requires MVV
  refusal cases for no match, ambiguous matches, and unmodeled recognized
  outcome.

Symptom: implementation ships tests for the pure table cases but not the
boundary cases. The first real RDR flow with a missing state read or an
unknown guard operator either guesses, panics, or collapses into a generic
internal error.

## 2. Section that will be rewritten within 6 weeks

`Normative Contracts` will be rewritten. It currently gives the correct
principle but not the reviewable result vocabulary. Once RDR 0005 maps CLI
codes and implementation adds tests, the project will need to retrofit the
kernel contract with the refusal kinds and result shape that should have been
locked here.

## 3. Assumption that will not survive first contact

A4 will not survive in its current form. It proves the CLI envelope is
extensible, but it does not prove that the kernel has a stable refusal
taxonomy for the CLI to map. `internal/cli/clierr::CLIError` can carry a code;
it does not decide which kernel refusal kinds exist or whether refusals are
ordinary values instead of errors.

## 4. Premortem

The resolver shipped with a clean-looking package and still broke the first
scripted skill integration. `Resolve` returned a transition plan on legal
inputs and an `error` on illegal inputs. The implementation had `ErrNoMatch`,
`ErrAmbiguous`, and `ErrUnknownOutcome`; missing owned state and unevaluable
guard both used `ErrInvalidInput` because the RDR's MVV did not force those
cases.

The `intrastate resolve --as=json` journey looked structured at first. It used
`internal/cli/respond::Fail` and emitted valid JSON. But the JSON code was the
CLI author's best guess, not a mapping from a kernel-owned refusal kind. RDR
flows saw `resolver-invalid-input` for both missing owned status and unknown
guard operator. Kata flows treated the same code as bad user input and retried
with a different recognized outcome, which could never fix an unavailable
owned snapshot.

Debugging started in the wrong layer. Engineers looked at Cobra flags and
artifact paths because the process exited through `CLIError`. The real defect
was that `Resolve` did not return a value-level `RefusalKind` that preserved
the exact kernel disposition. Tests passed because they covered legal replay,
no match, many matches, and unmodeled outcome. They did not cover unavailable
owned state or unevaluable guard, even though the normative contract named
both.

The repair required touching RDR 0001, RDR 0005, and the resolver package after
code existed. That is the failure: a planning document that explicitly existed
to lock the kernel contract left its central enum implicit until the CLI layer
had already invented dependent behavior.

## 5. Acceptance tests that would have caught the failures

Plain steps:

1. Build a fixture table with no matching edge.
2. Call `Resolve` with complete owned, observed, and recognized inputs.
3. Assert the returned disposition is a refusal value with kind `no_match`.
4. Assert the Go error return is nil or absent for the modeled refusal path.

Plain steps:

1. Build a fixture table with two matching edges.
2. Call `Resolve`.
3. Assert the returned refusal kind is `ambiguous_match`.
4. Assert the result carries enough source identity for RDR 0005 to map and
   report the ambiguous rows.

Plain steps:

1. Build a fixture table whose recognized-outcome alphabet excludes the
   supplied recognized outcome.
2. Call `Resolve`.
3. Assert the returned refusal kind is `unmodeled_outcome`.
4. Assert no persistence write description is present.

Plain steps:

1. Build a fixture table whose selected candidate requires an owned tag missing
   from the accessor-produced owned snapshot.
2. Call `Resolve`.
3. Assert the returned refusal kind is `owned_state_unavailable`.
4. Assert no ambient artifact discovery or accessor execution occurs.

Plain steps:

1. Build a fixture table whose candidate references a guard predicate that the
   guard evaluator reports as unevaluable.
2. Call `Resolve`.
3. Assert the returned refusal kind is `guard_unevaluable`.
4. Assert the CLI mapping layer can convert that kind to a stable
   `CLIError.Code` without inspecting an error string.
