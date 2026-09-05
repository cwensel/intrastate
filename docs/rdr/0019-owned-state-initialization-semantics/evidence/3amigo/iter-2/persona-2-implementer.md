Model: claude-opus-5[1m]

# Persona 2 — Implementer (iteration 2)

Re-run over the rewritten draft. Scope confined to the seven delta items in
the brief. Every code claim below was opened and read in
`/Users/cwensel/sandbox/newcoinc/intrastate`.

## Verdict

One Medium finding. Six of the seven delta items are sound and closed the
gaps the prior round opened. The seventh (item 4) fixed C1 and the
disposition table but left one Testing Strategy scenario asserting the
superseded claim — and the rewrite TOUCHED that scenario in the same pass,
so this is a new inconsistency, not a pre-existing one.

---

## Findings

### Medium — `0019:S5` now contradicts the rewritten `0019:C1`, and its expected refusal family is wrong against shipped code

**Blocks**: whether Phase 1 Step 2 owes a writer-arity fixture, and which
refusal family that fixture asserts.

C1's rewritten plan-validation clause says:

> the writer-arity half is already enforced at LOAD —
> `internal/table/load.go::checkAccessorBindings` folds every `[initial]`
> key into its `written` set and refuses `writerCount[key] != 1` as
> `CatMalformedAccessorBinding` — so no model reaching the verb can violate
> it ... The verb states the requirement ... but **owes no runtime check and
> no fixture for it.**

Grounded and TRUE. `internal/table/load.go::(*loader).checkAccessorBindings`
(line 1500) builds `written` from every row's `RequiresOwned` plus
`for _, t := range l.model.Initial { written[t.Key] = true }`, then refuses
`written[key] && writerCount[key] != 1` with
`"written tag %s is served by %d writers; want exactly one"`. Both the
zero-writer and multi-writer arms refuse at load. `checkAccessorBindings` is
step 103 in the loader's `run` chain, so no such model is constructible at
the verb.

The disposition table was correctly updated for this — the row changed from
`` `[initial]` key with ≠1 writer / role unbound `` to
`Role unbound (writer-arity is load-enforced, unreachable here)`. `0019:A3`
already agreed ("Because this runs at load, any model the CLI can see already
satisfies it by construction").

**`0019:S5` was not reconciled**, and it still reads:

> **Scenario**: `[initial]` names a key with zero declared writers, and
> separately one with more than one.
> **Expected**: refusal from the existing writer-routing family with ZERO
> writes committed — asserted on artifact bytes/mtime unchanged, or the
> artifact path still absent, not just exit code.

Two defects, both new relative to the rewrite:

1. **Unconstructible fixture.** By C1's own new reasoning the fixture model
   cannot load, so the verb never runs and there are no "artifact bytes" to
   assert on — the scenario's discriminating control has no subject. C1 says
   this half owes no fixture; S5 is exactly that fixture.

2. **Wrong refusal family.** A `checkAccessorBindings` failure is a LOAD
   failure, which the CLI surfaces as `codeModelInvalid = "flow-model-invalid"`
   (`internal/cli/flow_input.go:47`, emitted at `:203`; documented at
   `internal/cli/flow.go:134`, "A model that cannot be loaded refuses under
   flow-model-invalid"). The "existing writer-routing family" S5 names is
   `codeWriteUnbound = "flow-write-unbound"` / `::writerFor` — a different
   family and a different phase. An implementer building S5 from the record
   would write a test asserting `flow-write-unbound` and get
   `flow-model-invalid`.

This is not a stale leftover the rewrite merely failed to notice: the diff
shows S5's assertion clause WAS edited this pass (`asserted on artifact
bytes` → `asserted on artifact bytes/mtime unchanged, or the artifact path
still absent`), so the scenario was in hand while its premise was being
invalidated three sections above.

**Contagion**: `0019:S6` (role unbound — the half that IS reachable) now
reads "asserted on the same both-arms terms as S5", inheriting an
under-specified referent. S6 itself is sound — `internal/cli/flow_state.go:357`
shows the shipped pre-write role-binding loop (`if _, bound :=
req.artifacts[def.Accessor.Role]; !bound` → `codeArtifactMissing`), which is
genuinely reachable per invocation — but its assertion terms should not
point at a scenario the record elsewhere says is unbuildable.

**Cure** (one of): drop S5 and rewrite S6's assertion terms standalone; or
demote S5 to a LOAD-level scenario asserting `flow-model-invalid` at model
load with an explicit note that it is a loader regression guard, not a verb
test.

---

## Delta items verified clean — no finding

### 1. Kind-dispatched seed encoder (`0019:C1`) — CORRECT AND COMPLETE

`internal/cli/flow_input.go::canonicalValue` (line 714): for `decl.Kind != "set"`
the function checks `looksArray`, `value == ""`, and `table.ConformValue`, then
**returns `value` verbatim** — no coercion, no trimming, no normalization. So
`members[0]` is byte-identical to the scalar arm's output for every scalar
kind. `::canonicalSet` (line 770) is the set arm.

The brief's enumeration "enum, string, bool, int, float" does not match the
codebase: `internal/table/model.go:84` fixes `declaredKinds = {"enum", "bool",
"int", "set", "scalar"}` — there is no `string` kind and no `float` kind
(`IsDeclaredKind` refuses anything else at `tagDecl`, and `conformKind`/
`conformDomain` carry only `int`/`bool`/`enum`/`set` arms). C1 correctly says
"for every scalar kind", not an enumeration. **There is no kind where
`members[0]` diverges**, and the arity guard makes the index safe:
`::loadInitial` refuses `decl.Kind != "set" && len(members) != 1`, so a scalar
seed always has exactly one member. The A2 spike table confirms byte-identical
artifacts on all five scalar rows (enum/scalar-string/bool/int/float-on-scalar).

The only wrinkle is cosmetic, not a gap: `0019:A2` and `0019:S4` call the nine
spike rows "kinds" and list `float`, when `float` is a TOML *value shape*
authored on a `scalar`-kind tag (a float on an `int` tag is loader-refused —
the spike says so). That phrasing predates this rewrite and does not mislead
the encoder decision.

### 2. Carrier type-switch (`0019:C1` CARRIER SCOPE, `0019:A7`) — GROUNDED

All three claims hold:
- `internal/accessor/model.go:167` — `Binding Binding` is an exported field of
  exported `Definition`.
- `internal/cli/flowbind/registry.go:60-90` — the WRITE loop's switch yields
  exactly three types: default `&Writer{...}` (`*flowbind.Writer`),
  `case len(acc.Edit) != 0: NewEditWriter(acc, name)` (returns
  `&EditWriter{...}` per `edit.go:74-76`, so `*flowbind.EditWriter`), and
  `case commandBacked(acc): &cmdbind.Writer{...}`. **No fourth arm, no shared
  type between branches.**
- Exhaustiveness over production: `grep accessor.Definition{` finds exactly
  three construction sites, all in `registry.go`, and only the write loop
  yields write bindings. So `flowbind.Registry` is the single production
  writer of these types and A7's "To verify" item is confirmable.
- `commandBacked` is indeed unexported (`registry.go:121`), so C1's reason for
  type-switching instead of calling it is correct.

A7's statement matches what C1 says. A7 is still `Pending`; the fact is now
verifiable in one read, but flipping the status is the Resolve lens's call,
not a gap for me.

*One nuance worth an implementer note, not a finding*: `commandBacked` returns
true on `acc.Path == ""` (the carrier-less residue, `registry.go:122`), so a
carrier-less write accessor also yields `*cmdbind.Writer` and would refuse
under the carrier code while "naming the accessor and its carrier" — an
accessor that HAS no carrier. `0025:C1`'s load-time exactly-one rule makes
that residue unreachable through the loader, so the message never fires in
practice. No contract change needed.

### 3. ORDERING — carrier refusal preempts `--allow-commands` — IMPLEMENTABLE AS STATED

`flowbind.Registry(model, baseDir, allowCommands)` is called once, at
`internal/cli/flow_exec.go:1032`, during request assembly. It stores the gate
(`reg.AllowCommands = allowCommands`, `registry.go:43`) but **enforces nothing**
— construction is unconditional. Both allow-commands refusal sites fire strictly
later and only on invocation:
- `internal/cli/cmdbind/cmdbind.go:245` — `if !cfg.AllowCommands` inside
  `spawn`'s pre-spawn ladder.
- `internal/accessor/executor.go:388` — the edit-writer/command-reader
  precondition inside the write path.

So the constructed binding's dynamic type is knowable immediately after
registry construction and strictly before either gate. S10's requirement (run
the command-backed row with `--allow-commands` UNSET and expect the carrier
code) is achievable, and its "no helper binary is needed, since no accessor
runs" holds.

### 5. Non-normative refusal-code spellings — CONSISTENT

`flow-init-class-unsupported` and `flow-init-carrier-unsupported` appear in the
record at exactly two lines, both inside `0019:C1` (record lines 665-666), in
the same sentence that declares them non-normative. Every other site —
MVV step 8, `0019:S7`, `0019:S10`, the disposition table — says "the C1 class
code" / "the carrier code" / "class refusal code" / "carrier refusal",
referentially. **Nothing pins the literals.** The cited precedent is real:
`codeWriteEditRefused = "flow-write-edit-refused"`
(`internal/cli/flow_input.go:81`), and the `flow-<subject>-<condition>`
convention holds across the whole code block (`flow_input.go:36-107`).

### 6. Superset admission, arity-1 qualifier — CORRECT

`internal/table/load.go::loadInitial` (line 1680) carries verbatim:
`if decl.Kind != "set" && len(members) != 1 { return fail(CatMalformedInitialDeclaration, ...) }`.
So the array-literal-for-scalar admission is arity-1 only and a two-member row
never reaches the verb — exactly as C1 now says. The complementary direction
(bare scalar on a set tag) is unconstrained by that guard, which is why the
loader admits it, matching A2 spike row 10.

### 7. Downstream agreement with the rewritten C1

- `0019:A2` — now says "through the same kind-dispatched encoder `set-state`
  reaches (`::canonicalSet` for a set tag, the member verbatim for a scalar —
  C1)". Agrees.
- `0019:D-wire-byte-format` — makes no encoder claim; nothing to drift.
- `0019:D-selection-predicate` — carries the ALL quantifier, the STORE-key
  count including `::sealedKey`, the A6 carrier booking, and the file-backed
  scoping with the other two carriers refused. Agrees with C1.
- `0019:§mini-checks` `authority` row — "kind-dispatched: `::canonicalSet`
  (set) / `members[0]` (scalar)", with `::canonicalSet` for every kind now
  correctly listed as the SIBLING ARM (the rejected alternative). Agrees.
- `0019:§mini-checks` `trace` row — "encode + write | C1 canonical form,
  kind-dispatched ... | RT3 byte-identity vs `set-state`, both arms (A2 spike:
  9 kinds)". Agrees.
- `0019:§mini-checks` `disposition` row — updated to
  "Role unbound (writer-arity is load-enforced, unreachable here)". Agrees
  with C1; this is what makes S5's survival the anomaly.
- `0019:F5` (widened to — the empty-scalar disclosure sits on the encoder
  claim, so the encoder rewrite forced me here) — correctly rewritten from
  "decided wholly by what `::canonicalSet` renders for it" to "decided wholly
  by what C1's SCALAR encoder arm renders for it — `members[0]`, i.e. the
  empty string; `::canonicalSet` is the set arm and never reaches a scalar
  key." Grounded and right. F5 still says "C1's seed path assumes it cannot
  arrive" while `::loadInitial` demonstrably admits `note = ""` (no kind or
  domain rule rejects it, arity is 1) — but F5 discloses precisely that
  asymmetry and routes it to RDR 0002, and it is unchanged from the prior
  round. Not a new gap; not re-reported as a finding.

---

## Monday-morning residue (informational, no decision blocked)

With the Medium above resolved, the first hour is answerable from the record:
the encoder is a two-arm dispatch on `decl.Kind`; the carrier is a type-switch
on `Definition.Binding` admitting only `*flowbind.Writer`; the ordering is
carrier-then-gate and the code confirms it is possible; the writer-arity check
is not mine to write. The one thing I still cannot start is the emptiness
read, which `0019:A6` correctly books as undecided — that is a declared open
fork, not a gap.
