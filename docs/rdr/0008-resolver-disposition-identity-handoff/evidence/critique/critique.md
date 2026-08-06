Model: gpt-5.6-sol

# Critique — RDR 0008

Verdict: not lockable before the resolver trust boundary and the CLI consumer
obligation are made normative.

Pass status: first foundational-profile critique pass. The resources seam has
no alternate-model roster, so convergence still requires a fresh-session
single-model fallback pass and a diff; this file does not complete the lens.

## Origin ledger

### CRIT-1 — `Resolve` can accept a forged identity/action pair

**Root cause.** The draft trusted an upstream normalized-table check while the
production API remains `resolver.Input.Table []resolver.Edge`; callers and tests
construct exported rows directly. `matchingEdges` selects those values and
`plan` copies them, so a caller can pair rule B identity/locator with rule A
predicates/action and still return an authoritative-looking success.

**Enabling passages.** `Technical Design` assigned agreement to a boundary
before `plan`, and the original normative contract required only non-empty
identity fields at `plan`. Current code at
`internal/resolver/resolver.go::{Input,Edge,Resolve,matchingEdges,plan}` exposes
no validated or sealed table value. RDR 0007 plans to group outcomes and rows as
one semantic table, but grouping and revision hashing do not prove that authored
identity, locator, predicates, and action came from one source row.

**User symptom.** `flow resolve` applies the correct action but reports another
rule and source location. Replay and copy-isolation tests reproduce the same
false audit identity and therefore still pass.

### CRIT-2 — the locator-agreement validator is not implementable as written

**Root cause.** The draft fixed `SourceLocator` as an opaque string that `plan`
must not parse, while requiring an unnamed normalizer validator to prove that it
designates the same model/rule. RDR 0002 requires a locator to identify at least
model id and rule id, but does not pin a representation, parser, canonical
constructor, or comparison algorithm.

**Enabling passages.** `Technical Design` originally declared all four
`SelectedRule` fields as strings and simultaneously required semantic agreement.
Critical Assumption A4 was Pending and admitted that no callable constructor or
validator had been named.

**User symptom.** Implementers either reduce validation to presence, permitting
wrong diagnostic links, or invent an undocumented locator grammar that rejects
valid models when path or coordinate detail changes.

### CRIT-3 — the CLI can ship the anonymous plan this RDR claims to remove

**Root cause.** RDR 0008 deferred the only user-visible outcome to RDR 0005.
RDR 0005 `Technical Design` lists matched rule identity in its minimum payload,
but `Normative Contracts` currently require `flow resolve` to return only one
next tag-set or one refusal. An implementation following only MUST language can
omit selected identity and still pass.

**Enabling passages.** RDR 0008 `Implementation Plan / Phase 3` ended its own
acceptance at the carrier. RDR 0005
`Normative Contracts / flow resolve` does not repeat its non-normative payload
shape.

**User symptom.** `intrastate flow resolve --as=json` succeeds with next tags
but cannot identify the normalized rule that matched, reproducing the original
operator complaint.

## Section most likely to be rewritten within six weeks

`Technical Design` would be rewritten at the first normalizer integration. A
raw `SourceLocator string` plus upstream-trust prose cannot support the promised
agreement check. The implementation needs a typed locator or canonical
constructor and a production normalized-table value whose row components cannot
be fabricated independently.

## Assumption that fails first

A3 was incorrectly `Verified`. RDR 0005's descriptive payload includes matched
identity, but its normative `flow resolve` paragraph requires only a next tag-set
or refusal. The first real CLI consumer can therefore omit identity without
violating that RDR's current lock.

## Premortem

Six weeks after release, an operator investigates an unexpected transition.
`flow resolve --as=json` returns the expected next tags, so automation continues,
but the response has no enforceable matched-rule identity. A later text-mode
patch reconstructs a rule id from nearby normalized data and chooses a different
expansion of the same source rule.

The resolver tests are green. `Resolve` selects exactly one row; `plan` clones
the action; replay returns the same value twice. The production normalizer,
however, reused rule identity from the prior ordinary row while building an
escape row. Every identity field was non-empty, so the result was accepted. The
action belonged to the escape; the audit identity belonged to the ordinary
rule.

The promised normalizer check did not catch the mismatch because no constructor
or locator representation had been named. The locator was an opaque diagnostic
string, so the implementation treated agreement as a presence check. Replay
then fossilized the false pair and the mutation test proved only that the lie was
stable.

Ownership split the incident response: resolver maintainers pointed to RDR
0005, CLI maintainers pointed to its next-tag-only MUST, and normalizer
maintainers pointed to RDR 0008's unspecified locator. A CLI-side lookup restored
diagnostics quickly, recreating the exact second-authority path this RDR had
rejected. Operators stopped trusting `flow resolve` until the system was rebuilt
around a validated normalized row, typed source reference, and mandatory direct
projection.

## Acceptance tests working backward from failure

1. **CRIT-1 — forged composition.** Attempt ordinary and modeled-escape rows
   whose predicates/action come from rule A and whose identity/locator names rule
   B. Call the production resolver boundary. Construction must fail or the API
   must make the composition unrepresentable; no plan may be returned.
2. **CRIT-2 — locator agreement.** At the named normalizer
   constructor/validator, accept same-model/same-rule locator detail changes;
   reject another model, another rule, malformed input, and the zero value; copy
   the accepted typed locator into `TransitionPlan` unchanged.
3. **CRIT-3 — operator outcome.** For ordinary and write-free escape selections
   with identical actions but different identities, run production `flow
   resolve` in JSON and text modes. Require exact model id, rule id, expansion
   suffix, source locator, next tags, and writes from the plan, with no second
   lookup or accessor execution.
