Model: gpt-5.6-sol

# Critique — RDR 0008, fresh-context fallback pass

Fallback: single-model only; no alternate-model roster is configured. This pass
was run from a fresh context without reading the first critique or resolution.

## Origin ledger

### FALLBACK-1 — the selected escape has no defined action

**Root cause.** RDR 0008 makes an exactly-one modeled escape a successful
`TransitionPlan`, but RDR 0002 `Normative Contracts / escape` forbids an escape
rule from containing `write` or `clear`. RDR 0008 calls the escape write-free
without defining whether `NextTags` is empty, unchanged, derived, or authored.

**Enabling passages.** `Proposed Solution / Approach`, `Technical Design`, and
Validation scenario 3 require the same plan carrier and “expected next tags.”
RDR 0002's `draft-no-match-escape` example has no action. Current
`internal/resolver/resolver_test.go::TestResolve_AmbiguityRequiresOneExplicitEscapeEdge`
hand-constructs an escape `Edge` with next tags and writes that the source
schema rejects.

**User symptom.** `flow resolve` reports success for an escape but returns an
empty or invented action, so automation repeats the same state or executes a
transition the authored model never declared.

### FALLBACK-2 — selected-rule identity is not historical event identity

**Root cause.** `SelectedRule` contains `(model id, rule id, expansion suffix)`
and diagnostic source location but not `TableRevision`. The tuple identifies a
logical row only within one normalized table, while the problem statement can
be read as promising a durable audit identity across later model revisions.

**Enabling passages.** `Load-Bearing Decisions / Identity` defines tuple
equality and excludes source locator from logical identity. RDR 0007 `Normative
Contracts` defines revision over the complete normalized semantics and its
`Existing Infrastructure Audit` charts enforcement of the table/revision
association to a successor.

**User symptom.** An archived selection and a later, semantically changed table
reuse the logical tuple; an operator follows current semantics while diagnosing
an older decision.

### FALLBACK-3 — validated-table immutability is asserted, not proved

**Root cause.** The draft requires callers to be unable to compose or replace
row parts, but A4/A5 do not yet name the production type or constructor. The A2
spike uses caller-assembled structs, a string locator, and mutations only after
`plan` returns.

**Enabling passages.** `Technical Design` promises an immutable validated table;
`internal/resolver/resolver.go::{Input,Edge}` exposes public row storage;
`evidence/spikes/ownership_test.go::TestPlanOwnsSelectedIdentityAndAction` does
not mutate constructor inputs before `Resolve`, exercise a typed locator, or
prove an external-package API boundary.

**User symptom.** A caller or reload mutates aliased row storage between
validation and `plan`, producing a coherent-looking selection whose identity
and action came from different snapshots.

## Section rewritten within six weeks

`Technical Design` would be rewritten first. The concrete normalized-table API,
escape disposition, typed-locator ownership, and revision scope all determine
what `TransitionPlan` can truthfully promise; field-copy prose cannot substitute
for those boundaries.

## Assumption that fails first

A2 fails first. Its Verified evidence proves a string locator and post-plan
map/slice isolation, while the decided design now requires a typed locator and
a snapshot-owned validated table. It does not prove the current claim.

## Premortem

The parser accepted an authored no-match escape with no writes, as RDR 0002
requires. `resolver.Resolve` selected it and `plan` returned a successful plan
with an empty action. The CLI emitted `ok`; the workflow retried against the
unchanged state forever. Separately, a cached table reload mutated backing
storage after validation but before `plan` cloned the action, so a diagnostic
named one rule while applying another row's tags. The replay test stayed green
because it called `Resolve` twice over one in-memory input. An archived record
later pointed at a changed logical rule because readers mistook `SelectedRule`
for revision-bound historical identity. Repair required defining the escape
disposition, making the normalized table an owned snapshot, and making revision
scope explicit.

## Acceptance tests working backward

1. Parse valid RDR 0002 TOML containing an escape with no write/clear block,
   normalize it, resolve it, and render both CLI modes. Require the one explicit
   disposition/action meaning chosen across all four boundaries; never
   hand-construct an `Edge`.
2. Construct a validated table from caller-owned row slices, maps, writes, and
   locator inputs; mutate all inputs before and after `Resolve`. External-package
   API tests must show row parts cannot be replaced, and `go test -race` over a
   stable snapshot must remain coherent.
3. Resolve two normalized tables that share the logical selected-rule tuple but
   have different RDR 0007 revisions. The selection must be documented as
   revision-scoped, and the revision-binding successor—not this carrier—must
   prevent a free label from being paired with different semantics.
4. Load malformed authored identity in JSON and text modes. Require a stable
   load/config `CLIError`; reserve the internal programmer error for an
   impossible invalid value inside an already validated table.
