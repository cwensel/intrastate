# Critique - Guard Predicate Exhaustiveness

## 1. Three Ways Implementation Goes Wrong

### Failure 1: Lint proves the wrong shape of exhaustiveness

Root cause: the RDR says coverage is checked over "the checked
dimension/product" but never names who defines that product or how rows are
grouped before coverage/overlap are computed.

Enabling passage: Critical Assumption A2 says coverage is
`union(row_i accepted assignments) == D` for "the checked dimension/product";
Technical Design repeats that static lint expands atoms into finite domain
constraints but does not bind the row group or product shape.

User symptom: the first lint implementation checks one tag dimension at a time,
or checks all rows globally, and misses a gap that only exists for a specific
state/outcome/profile combination. A flow author sees `intrastate lint` pass,
then the resolver returns `no_match` or `ambiguous_match` for a legal transition.

### Failure 2: Set-valued domains become a hidden exponential cost

Root cause: the RDR permits set-valued tags by treating the declared element
universe as a finite powerset, but it does not constrain the implementation to a
symbolic/bitset representation or tell lint what to do when the universe is
large enough to make naive enumeration unreasonable.

Enabling passage: Critical Assumption A2 allows "set-valued tags as the
declared element universe's finite powerset or an implementation-equivalent
bitset"; Technical Design says set-element universes are finite by declaration;
Performance Expectations only says local typed comparisons and finite-set
expansion, without a refusal/downgrade behavior for oversized finite domains.

User symptom: a seemingly valid transition model with 15 declared lens or tag
elements hangs lint, allocates heavily, or gets an implementation-specific cap
that is not in the contract. The user cannot tell whether the model is invalid
or the tool is slow.

### Failure 3: Predicate error ownership leaks across RDRs

Root cause: this RDR says predicate parse and lint failures can add stable CLI
codes, while RDR 0005 and RDR 0006 own the user-facing CLI command and graph-lint
diagnostic surfaces. The semantic error kind, CLI code, and lint finding code can
drift because ownership is split but not named.

Enabling passage: A4 and Existing Infrastructure Audit say predicate parse,
unknown-operator, type-mismatch, non-exhaustive, overlap, and unevaluable-guard
failures can add stable codes on the existing gateway. Capability Dependencies
list RDR 0005 as pending for CLI output and RDR 0006 as pending for graph lint,
but the Proposed Solution does not define the handoff.

User symptom: predicate parse errors, graph lint findings, and resolver runtime
refusals use different names for the same condition. Text mode says one thing,
JSON mode emits another code, and implementation tests end up asserting wrapper
behavior instead of predicate semantics.

## 2. Section Rewritten Within Six Weeks

The `Technical Design` section will be rewritten. It currently has the right
operator matrix, but the finite-domain proof is underspecified at the exact
point implementers need it: grouping candidate rows, building the guard product,
handling set-valued dimensions without materializing an exponential powerset,
and separating predicate semantic errors from CLI/lint envelopes.

## 3. Assumption That Will Not Survive First Contact

A2 will not survive unchanged. "Every exhaustiveness claim can be reduced to
finite declared domains" is true mathematically but incomplete operationally:
real lint needs a scoped product, a representation strategy, and a refusal or
downgrade path for finite-but-too-large domains.

## 4. Premortem

We shipped guard predicate exhaustiveness and the first serious user journey
failed at the boundary between `all`/`unless` authoring and graph lint. A flow
author encoded the RDR prelock path: profile routes to different lens sets,
`prelock_iterations.gte = 3` disables another pass, and set-valued lens tags
drive whether `critique` or `repeatability` is still required.

The implementation normalized rows per RDR 0002 and fed them to the RDR 0006
lint. The predicate package correctly parsed `eq`, `in`, `gte`, `exists`, and
`contains`; `resolve` used the same atoms at runtime. The bug was that
exhaustiveness was checked per tag dimension, not per selection group and
multi-tag product. The table passed because each individual dimension looked
covered. At runtime, `resolve` saw a legal tag-set whose profile/lens/iteration
combination matched no row, and returned `no_match`.

The second failure came when we fixed that and added set-valued products
literally. A user declared a larger set universe than the fixture had modeled.
The lint command became slow and memory-heavy. Because the RDR did not specify a
symbolic representation or a too-large-domain diagnostic, one implementation
silently capped enumeration while another kept expanding. The same model passed
in one environment and failed in another.

The third failure was diagnostic drift. `internal/cli/clierr.CLIError` and
`internal/cli/respond.Fail` were available, but the predicate package, graph lint
package, and future CLI command each coined names. Users saw
`unknown-operator`, `guard_operator_unknown`, and `lint-predicate-operator` for
the same authored mistake. Tests hardened around accidental strings instead of
the semantic contract.

## 5. Acceptance Tests That Would Catch This At RDR Review Time

1. Given a normalized row group for one source state and recognized outcome with
   two varying guard dimensions, when lint checks coverage, then it evaluates
   the Cartesian product of the participating finite domains for that row group
   and detects a gap that is invisible in either single dimension alone.

2. Given a set-valued tag with a declared element universe and a `contains`
   predicate, when lint proves coverage/overlap, then it uses a deterministic
   symbolic or bitset-equivalent representation and never requires source-order
   priority or host callbacks.

3. Given a finite declared domain whose product is too large for the configured
   lint implementation to prove, when lint evaluates exhaustiveness, then it
   refuses or downgrades that exhaustiveness claim with a stable diagnostic
   instead of silently capping enumeration.

4. Given an unknown predicate operator in authored source, when the model is
   parsed/linted and surfaced through the CLI, then the predicate semantic kind,
   lint finding code, and CLI envelope have one documented ownership chain and
   preserve the source rule/context id.
