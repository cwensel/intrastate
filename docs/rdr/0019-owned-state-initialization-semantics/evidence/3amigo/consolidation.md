Model: claude-opus-5[1m]

# 3amigo consolidation — 0019

Three isolated persona passes (PM, Implementer, QA), no cross-persona
visibility. Hotspots below are a mechanical count over `rdr anchors` ids, not a
re-judgement: an id raised by ≥2 personas marks a **hotspot passage**, not a
validated finding, and a single-persona finding is not thereby weaker.

## Hotspots (id cited by ≥2 personas)

| id | personas | what converged |
|---|---|---|
| `0019:C1` | PM, Impl, QA | 3-way. The contract every persona had to widen into: carrier scope (PM: unstated in the promise), named encoder + carrier discriminator (Impl), ALL quantifier untested (QA). |
| `0019:MVV` | PM, Impl, QA | 3-way. Fixture-only validation (PM); reads payload field names that C1 defers (Impl); step 8 asserts an unspelled code (QA). |
| `0019:A6` | PM, Impl | Pending emptiness carrier — the predicate everything rests on. |
| `0019:A5` | Impl, QA | The resurrection hazard the ALL quantifier exists to exclude. |
| `0019:S4` | Impl, QA | Loader-only kinds: fixture arity (Impl) vs RT3 disagreement (QA). |
| `0019:S7` | Impl, QA | Refusal-code spelling absent from the oracle. |
| `0019:S8` | Impl, QA | Read-back mismatch unconstructible on the admitted carrier. |
| `0019:S9` | Impl, QA | Sealed-store fixture. |
| `0019:S10` | Impl, QA | Carrier refusal: unexported discriminator (Impl), ordering vs `--allow-commands` (Impl), no code spelling (QA). |
| `0019:D-selection-predicate` | Impl, QA | The predicate's encoder arm and its untested quantifier. |
| `0019:§approach` | PM, Impl | Reads unqualified/singular where C1 is scoped/plural. |
| `0019:§technical-design` | Impl, QA | Deferred spellings and field names. |
| `0019:§mini-checks` | Impl, QA | Authority/oracle/disposition tables: blank A6 cell; no quantifier row. |

## Merged findings, severity-ranked

### High

- **H1 `0019:§background` (PM)** — the Priority: Low justification is factually
  false of the file it names. `models/rdr.toml` declares three owned tags with
  `provenance = "owned"`, an `[initial]` block seeding all three, no `class` key
  (so state-machine per `internal/table/model.go::IsDecisionTable`), and a
  file-backed `[write.rdr-status]`. It is squarely in the population the verb
  serves. Blocks: the Priority field and the force of `§decision-rationale`'s
  rejection of Alternative 1.
- **H2 `0019:C1` / `0019:D-selection-predicate` (Impl)** — the named seed
  encoder has the wrong signature for scalar `[initial]` values.
  `internal/cli/flow_input.go::canonicalSet` is `func(members []string) string`,
  called only from `::canonicalValue`'s `isSet` arm; the scalar arm returns the
  value verbatim. Six of `0019:S4`'s nine kinds are scalars, and no element
  names an encoder that renders them correctly. Blocks: the
  `Model.Initial → []resolve.Tag` function and whether it dispatches on kind.
- **H3 `0019:C1` (carrier scope) / `0019:S10` (Impl)** — the discriminator the
  contract keys on is unreachable. `internal/cli/flowbind/registry.go::commandBacked`
  is unexported; the verb lives in `internal/cli`. Two implementable readings
  survive (type-switch on `accessor.Definition.Binding` for `*flowbind.Writer`;
  or re-derive from `table.Accessor`, which `D-selection-predicate` bans) and
  the record picks neither. Blocks: where the carrier check lives and what it
  reads; `S10`'s negative control.
- **H4 `0019:A6` (Impl, PM)** — Pending, and it is the carrier for the predicate
  C1 and `D-selection-predicate` are entirely built on. Persona 2 independently
  re-grounded and confirmed the record's gap analysis is correct. Blocks: the
  emptiness probe's signature and owning package; `MVV` step 2 and `S9` cannot
  be attempted before it lands. A6's own "If wrong" calls a refutation
  approach-level.
- **H5 `0019:S9` (QA)** — the sealed-store fixture is not constructible as
  written. The seal keys on the writer's declared locator
  (`flowbind::Writer.Apply` under `if unreachable(w.Path)`, a pure suffix test),
  so a writer is always- or never-sealing; `0019:A3` forces the clears through
  that same writer, which then exits non-zero via `read_back_incomplete`.
  Prevents: the only test pinning the emptiness count to STORE keys rather than
  owned keys — its negative control's whole point.
- **H6 `0019:C1` / `0019:D-selection-predicate` (QA)** — the ALL quantifier,
  declared "load-bearing rather than stylistic", has no scenario, MVV step, or
  mini-check row. Every fixture is single-artifact. Prevents: bind two roles,
  seed A, leave B empty, assert no-op over both — the change that resurrects
  into a still-keyed store, which is the `0019:A5` hazard.

### Medium

- **M1 `0019:§problem-statement` + `0019:§approach` (PM)** — the stated user
  outcome is unqualified ("every state-machine flow") where `C1` and
  `§consequences` establish a file-backed-only carve-out. Blocks: what Phase 2's
  reference docs and `--help-all` may claim.
- **M2 `0019:MVV` (PM)** — validation never runs against a model that ships;
  every step is fixture-scoped. Blocks: whether the MVV suffices for
  `0019:G-scope`.
- **M3 `0019:C1` / `§mini-checks` disposition table (Impl)** — the writer-arity
  refusal arm is unreachable for any admitted model:
  `internal/table/load.go::checkAccessorBindings` already folds `[initial]` keys
  into the `written` set and refuses `writerCount[key] != 1` at LOAD. The
  role-unbound half IS reachable. Same UNCONSTRUCTIBLE argument C1 makes for the
  decision-table class, but here a live refusal arm is kept.
- **M4 `0019:C1` vs `0019:S10` (Impl)** — refusal ordering against a
  command-backed writer with `--allow-commands` unset is unnamed; two exit-2
  refusals are available and the record does not say which fires.
- **M5 `0019:S8` (QA)** — no constructible read-back mismatch on the only
  carrier C1 admits: writer and reader share one store. The constructions that
  would work (reader bound to a different path; command-backed *reader*, legal
  since the carrier refusal is scoped to WRITE accessors) are unstated. Prevents
  the PRESENT-AND-UNVERIFIED arm and its anti-repair clause.
- **M6 `0019:F2` (QA)** — the torn multi-writer failure mode has no scenario;
  the disposition table collapses it into the generic non-empty row, and on a
  single-writer fixture torn and post-clear are indistinguishable. Prevents the
  test that fails when an implementation "helpfully" completes a torn seed.
- **M7 `0019:C1` / `0019:S7` / `0019:S10` / `MVV` step 8 (Impl, QA — both
  personas, independently)** — the two new refusal codes have no spelling and
  the payload has no field names, so the oracles cannot be written as
  equalities. Disclosed deferral ("sharpened pre-lock"), i.e. owed now.

### Low

- **L1 `0019:§consequences` (PM)** — the two-command first-run cost is recorded
  but no passage states the verb is typed by hand, and no discovery path points
  an operator at it. `BR4`/`BR5` close both automatic routes.
- **L2 `0019:§approach` (Impl)** — reads singular ("the bound artifact is
  EMPTY") where C1 is emphatically plural. C1 governs; named because Approach is
  what an implementer skims first.
- **L3 `0019:C1` (Impl)** — "an array literal for a scalar tag" loads only at
  arity 1: `::loadInitial` refuses `decl.Kind != "set" && len(members) != 1`.
  The claim is true as narrowed; the qualifier is absent.
- **L4 `0019:RT3` vs `0019:S4` (QA)** — RT3 and S4 state different things about
  the same three loader-only kinds. The `fidelity` mini-check row already
  reconciles them; RT3 is where a test author would stop reading.
- **L5 `0019:S5`, `0019:S6` (QA)** — "asserted on artifact bytes" is undefined
  when no artifact exists; `S7` two scenarios earlier states both arms.

## Non-findings recorded by the personas

- The problem statement names a user-visible outcome, not an engine irritation:
  `internal/graphlint/analysis.go` emits "the model declares no initial owned
  state" while no non-test `internal/cli` code reads `Model.Initial` — the
  lint-demands / runtime-ignores gap is real (PM, grounded).
- Scope did not expand past the problem; `BR1`–`BR5` record five adjacent
  temptations as rejected (PM).
- Every other `::Symbol` C1 and the mini-checks cite resolves with the shape the
  record assumes (Impl enumerated them). One attribution correction:
  `::checkClassAgreement` is in `internal/table/load.go`, not `graphlint`.
- `G-scope` / `G-proportionality` are unanswered template placeholders — Stage
  7's to answer, not a pre-lock finding (PM).
