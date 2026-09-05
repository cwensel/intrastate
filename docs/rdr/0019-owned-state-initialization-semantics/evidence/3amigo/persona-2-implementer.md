Model: claude-opus-5[1m]

# Persona 2 — Implementer

Question asked of the record: if I started coding `flow init-state` Monday
morning, what would I ask in the first hour?

Starting set: `0019:C1`, the four `D-*` elements, and the `source-anchor`
edges. **Widening declared**: every code claim below was ground-checked
against the source repo, and three findings are about what the contracts do
NOT say (a scalar-encoding arm, a carrier-detection call site, a refusal
ordering) — silences have no line range, so I read `0019:§technical-design`,
`0019:§mini-checks`, `0019:§failure-modes`, `0019:A5`, `0019:A6`,
`0019:S4/S7/S9/S10`, and `0019:MVV` to confirm the silence was real and not
answered elsewhere in the record. What sent me: C1's `::canonicalSet`
citation did not typecheck against `Model.Initial`'s shape, and C1's
`::commandBacked` citation named an unexported symbol.

---

## High

### 1. `0019:C1` / `0019:D-selection-predicate` — the named seed encoder has the wrong signature for scalar `[initial]` values; the scalar arm is unwritten.

C1 fixes the seed rendering on one encoder: values are "rendered in JDR 0001
§D13's canonical form (the same encoder `set-state` uses,
`internal/cli/flow_input.go::canonicalSet`)". The `mini-checks` `authority`
table repeats it — "Value encoding | `::canonicalSet` … the byte-equality
source" — and `0019:§failure-modes` F5 leans on it again ("decided wholly by
what `::canonicalSet` renders for it").

Grounded: `internal/cli/flow_input.go::canonicalSet` is
`func canonicalSet(members []string) string` and has exactly ONE production
call site — `internal/cli/flow_input.go::canonicalValue`, inside its
`isSet` arm (`return canonicalSet(members), nil`). The scalar arm of
`canonicalValue` returns `value` verbatim and never touches `canonicalSet`.
Meanwhile the seed source is `internal/table/model.go::TagValue`
(`{Key string; Value []string}`), and the write seam is
`resolve.Tag{Key, Value string}`.

So for a scalar `[initial]` key the verb has `[]string{"draft"}` in hand and
the record's named encoder would render it `["draft"]` — a set literal into a
scalar tag, which is not what `set-state --write status=draft` persists.
The correct rendering is `members[0]`, and no element says so. `set-state`'s
actual encoder is `canonicalValue`, but C1 explicitly rules that route out
(the argv route "would refuse at seed time models that load clean"), so the
implementer cannot simply call it either.

**Decision it blocks**: what function to write for
`Model.Initial → []resolve.Tag`, and whether it dispatches on
`m.Tags[key].Kind == "set"`. This is the exact site RT3 (byte-identity with
the `set-state --write` transcription) and `0019:S4`'s nine-kind table
assert against — six of those nine kinds are scalars, and the record names no
encoder that produces a correct value for any of them.

### 2. `0019:C1` (carrier scope) / `0019:S10` — the carrier discriminator the contract keys on is unexported and unreachable from the verb's package.

C1's CARRIER SCOPE clause requires init-state to refuse "against a model any
of whose bound write accessors is non-file-backed", and identifies the
branch: "`::Registry` selects among THREE … branching on `len(acc.Edit) != 0`
to the line-oriented `::NewEditWriter` (0028) and on `::commandBacked` to
`cmdbind.Writer` (0025)". `0019:S10` restates it — "routed to
`flowbind::commandBacked`".

Grounded: `internal/cli/flowbind/registry.go::commandBacked` is
**unexported** (`func commandBacked(acc table.Accessor) bool`). The verb
lives in `internal/cli`, a different package. The only exported names in
`flowbind` are `Registry`, `OwnedTags`, `Reader`, `Writer`, `Gate`,
`EditWriter`, `NewEditWriter`. So the verb cannot call the discriminator the
contract names.

Two implementable readings survive, and the record picks neither:
(a) type-switch on `accessor.Definition.Binding` for `*flowbind.Writer`
(`internal/accessor/model.go::Definition.Binding` is an exported interface
field, `*flowbind.Writer` is exported — this works); or
(b) re-derive from `Definition.Accessor` (`table.Accessor`'s `Edit`,
`Command`, `Path`) in `internal/cli`, duplicating `commandBacked`'s body —
which is precisely the "re-derivation" `0019:D-selection-predicate` bans for
the class question ("never a re-derivation from `len(owned)`").
A third option — export `commandBacked` — is a `flowbind` API change the
record does not authorize.

**Decision it blocks**: where the carrier check lives and what it reads.
`0019:S10`'s negative control ("removing the carrier check makes the edit
case attempt a write") cannot be written until this is settled.

### 3. `0019:A6` is Pending, and it is the carrier for the predicate `0019:C1` and `0019:D-selection-predicate` are entirely built on.

C1 states outright that "the emptiness read is a SECOND new data flow" and
that "Which carrier serves it … is UNDECIDED here and is booked as A6".
`0019:A6` Status: Pending. The `mini-checks` `authority` table has the cell
blank — "Store emptiness | the artifact's own key presence | **no shipped
reader** — carrier undecided (A6)" — and the `trace` table marks the
emptiness step "**A6 carrier undecided**".

Grounded and confirmed correct: `internal/accessor/binding.go::ReadBinding`
is `Read(ctx, art, requested []string) ([]KeyValue, []string, error)` —
key-scoped, no cardinality. `internal/cli/flowbind/flowbind.go::store` is
`type store map[string]string`, package-private, and no exported function
returns a count. `0004:C3` forbids reading the artifact directly. The record
is right that nothing shipped can answer this.

**Decision it blocks**: the single largest one — the signature of the
emptiness probe, which package owns it, and whether it changes
`accessor.ReadBinding` (which would touch the other implementer,
`internal/cli/cmdbind/cmdbind.go::Reader.Read`). `0019:MVV` step 2 and
`0019:S9` cannot be attempted before this lands. A6's own "If wrong" says a
refutation is "approach-level, not editorial", so this is not a detail to
discover during implementation.

---

## Medium

### 4. `0019:C1` (plan validation) / `0019:§mini-checks` disposition table — the writer-arity refusal arm is unreachable for any admitted model.

C1: "every `[initial]` key must route to exactly one declared writer and
every needed artifact role must be bound, or the verb refuses with ZERO
writes committed". The disposition table gives it a row: "`[initial]` key
with ≠1 writer / role unbound | 2 | zero (plan-level)".

Grounded: `internal/table/load.go::checkAccessorBindings` (`run`'s last step,
`load.go:103`) already enforces this at LOAD. Its comment is explicit —
"Every key any rule's write block or clear list names, and every key
`[initial]` assigns, MUST be served by exactly one writer" — and the loop
`for _, t := range l.model.Initial { written[t.Key] = true }` folds
`[initial]` keys into the `written` set, whose arm refuses
`CatMalformedAccessorBinding` when `writerCount[key] != 1`. Any model that
reaches the verb has already passed this.

This is the same UNCONSTRUCTIBLE argument C1 spends a paragraph making for
the decision-table class ("no admitted decision-table model carries
`[initial]`") — but here the record keeps a live refusal arm instead. The
role-unbound half of the row IS reachable (`runFlowSetState` checks
`req.artifacts[def.Accessor.Role]` at runtime, `flow_state.go:357-363`).

**Decision it blocks**: whether to write a writer-arity check in the verb
(and a test for it that cannot be given a fixture), or to call
`internal/cli/flow_state.go::writerFor` only for its role-binding side
effect. Also whether `0019:S8`'s "unbound artifact role" case is the only
testable half of that MVV row.

### 5. `0019:C1` (carrier scope) vs `0019:S10` — refusal ordering against a command-backed writer with `--allow-commands` unset is unnamed.

C1 requires the carrier refusal to be "a distinct terminal refusal in the
`flow-*` family naming the accessor and its carrier". `0019:S10` runs the
command-backed case and expects "exit 2 with the carrier refusal code …
with ZERO writes and no accessor invocation", noting over-determination only
for the EDIT carrier.

Grounded: a command-backed accessor is also gated by `--allow-commands`
(`internal/accessor/model.go::Registry.AllowCommands`,
`internal/cli/cmdbind/cmdbind.go` `if !cfg.AllowCommands`), and
`internal/accessor/executor.go` already carries an exit-2 refusal for it
("requires the allow-commands opt-in (--allow-commands); refusing before
mutation rather than …"). So the command-backed row of S10 has TWO exit-2
refusals available and the record does not say which fires. If the fixture
omits `--allow-commands` the test may pass on the wrong code; if it sets it,
the fixture needs a real helper binary.

**Decision it blocks**: whether the carrier check runs before registry
construction / executor entry (preempting the allow-commands refusal), and
which flag the `0019:S10` command-backed fixture sets.

### 6. `0019:C1` — the two new refusal codes have no spelling, so the payload/refusal constants cannot be written.

C1: "It carries a dedicated code in the `flow-*` family (spelling sharpened
pre-lock)" for the class refusal, and "a distinct terminal refusal in the
`flow-*` family" for the carrier refusal. `0019:D-wire-byte-format` defers
payload field names too: "Field names deferred to Resolve/Pre-Lock, owned
here."

Grounded: `internal/cli/flow_input.go:36-107` is the code constant block
(`codeTagInvalid = "flow-tag-invalid"` … `codeArgvUpstreamStop`); there is
no `init` code and no `init-state` command anywhere in `internal/`
(`internal/cli/flow.go:95-98` registers `next`, `resolve`, `read-state`,
`set-state` only). Two new constants and one new payload struct must be
minted with no name given.

**Decision it blocks**: the constant spellings that `0019:S7`, `0019:S8` and
`0019:S10` assert on, and the payload field names that `0019:MVV` steps 2 and
7 read ("the seeded-all case" vs "which `[initial]` keys the store does not
carry" — C1 fixes the semantics but not one field name).

---

## Low

### 7. `0019:§approach` reads singular where `0019:C1` and `0019:D-selection-predicate` are emphatically plural.

Approach: "init seeds if and only if the bound artifact is EMPTY". C1: "the
quantifier is ALL, not ANY and not per-artifact, and it is load-bearing
rather than stylistic". `D-selection-predicate` agrees with C1 ("EVERY bound
artifact carries no key").

Grounded that the multi-artifact case is real: `internal/table/model.go`
`Model.Writers map[string]Accessor`, and
`internal/cli/flow_state.go::runFlowSetState` iterates
`slices.Sorted(maps(byWriter))` checking `req.artifacts[def.Accessor.Role]`
per writer — N writers over N roles is constructible.

C1 governs and is unambiguous, so this is a low-severity second reading in
prose rather than a live fork.

**Decision it blocks**: nothing, if the implementer reads C1 first. Named
because Approach is the section an implementer skims before the contract.

### 8. `0019:C1` — "an array literal for a scalar tag" loads only at arity 1.

C1 lists the loader's superset admissions: "a bare scalar for a set-valued
tag and an array literal for a scalar tag both load and normalize".
`0019:S4` promotes this to a test row ("array for a scalar tag").

Grounded: `internal/table/load.go::loadInitial` refuses
`decl.Kind != "set" && len(members) != 1` with
`CatMalformedInitialDeclaration` ("kind %s holds one value, not a member
sequence"). So `["draft"]` for a scalar loads; `["draft","final"]` does not.
The claim is true as narrowed but the qualifier is absent, and an
implementer building the `0019:S4` fixture from C1's wording alone would
write a two-member row that refuses at load.

The related `checkClassAgreement` citation in C1 resolves — it is
`internal/table/load.go::checkClassAgreement` (`load.go:819`, message
`"[model] class %s declares owned=%d; a %s model declares …"`), not a
`graphlint` symbol as the surrounding `::checkDanglingEdge` /
`::checkAlwaysPresentOwned` citations might suggest. `::loadInitial`,
`::IsDecisionTable`, `::Model.Initial`, `::sealedKey`, `::NewEditWriter`,
`::Writer.Apply`, `::load`, `::writerFor`, `::groupByWriter`,
`::parseWrites`, `::canonicalValue`, `::ReadBinding.Read`, `::store`,
`::KeyValue.Absent`, `::buildRequest`, `::runFlowSetState`, `::reach`,
`::checkAlwaysPresentOwned`, `::checkDanglingEdge` all exist with the shape
the record assumes.

**Decision it blocks**: the `0019:S4` "2 kinds only the loader admits"
fixture arity.
