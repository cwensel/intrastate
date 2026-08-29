Model: claude-opus-5[1m]

# Persona 2 — Implementer

Question answered: *if I started coding this Monday, what would I ask in the first hour?*

Owned starting set: `elements[]` kind `C` (`0023:C1`, `0023:C2`), kind `D`
(`0023:D-identity`, `0023:D-wire-byte-format`, `0023:D-naming`,
`0023:D-selection-predicate`), and the `source-anchor` edges.

**Widened.** Four times, each named at the finding. What sent me out of the
owned set was, in every case, a *silence* in C1/C2 that has no line range:
C1 mandates a total command-tree walk but never says how the tree is
materialized (→ `§implementation-plan` / `0023:A5` / cobra source behaviour);
C2 mandates a reflective partition oracle but names no in-code home for the
assignment (→ `§source-authority-census`); Phase 3 edits a help string that a
committed generated file is pinned against (→ `§implementation-plan` Phase 3 +
Makefile); and the Prerequisites checkbox contradicts a projected element
status (→ `§prerequisites` vs `0023:A5`).

---

## S1 — Blocking

### F1. `0023:C1` (widened: `0023:A5`, `0023:§prerequisites`, `0023:§implementation-plan` Phase 2) — the whole-tree walk cannot see `completion` on the tree C1 hands it, and C1 does not say who materializes it

C1's clause is emphatic — *"WHOLE means total: the walk descends every child
of the root with no name-based skip and no Hidden gate, **INCLUDING the
auto-generated help and completion commands**"* — and `§oracle-discriminability`
S4 adds a negative control asserting the walk actually reaches `help`.

But cobra does not create `completion` at `AddCommand` time. I probed the real
root:

```
NewRootCmd().Commands()                      → docs, flow, help, lint, version
+ InitDefaultHelpCmd(), InitDefaultCompletionCmd() → completion, docs, flow, help, lint, version
```

`help` exists on a fresh `NewRootCmd()` only because `registerHelpAllOnTree`
already forces it (`internal/cli/help_all.go::wireHelpSubcommandAll`).
`completion` appears only after `InitDefaultCompletionCmd()`, which cobra runs
inside `ExecuteC` — never in a unit test that builds the tree and walks it, the
idiom every shipped absence oracle uses
(`internal/cli/flow_all_0011_test.go:460`, `internal/cli/flow_input_0005_test.go:460`).

So the S4 oracle as specced passes **vacuously on `completion`**: it "includes"
a command that is not there. The four-verb vacuity guard C1 inherits from
`TestReq46And47And65And98_…` does not catch this — it guards the *flow group*,
not the root's auto-generated children. C1's own justification for the total
walk is the `help --all` precedent (a registration reachable only on an
auto-generated command), which is exactly the class this gap re-admits.

**Clarification requested:** does the S4 oracle call
`root.InitDefaultCompletionCmd()` (and `InitDefaultHelpCmd()`) before walking,
or does C1's "INCLUDING … completion" only bind whatever the tree happens to
carry? If the former, C1 or S4 must say so — it is not derivable from the text
and it is the difference between a real oracle and a green vacuity.

**Blocks:** Phase 2's whole-tree registration oracle — the single piece of new
machinery C1 mandates and A5 says has *no shipped exemplar*. I cannot write it
without knowing whether the negative control (`help` is reached) is the whole
guard or whether `completion` needs its own.

**Bonus ask on the same clause:** should the vacuity guard be strengthened to
assert the walk reaches *both* `help` and `completion`? As written, S4's control
names `help` only — the one that happens to already exist.

### F2. `0023:§prerequisites` vs `0023:A5` — the Prerequisites checkbox is unsatisfiable on the record as projected

`§prerequisites` reads `- [ ] All Critical Assumptions verified`. The projector
reports `0023:A5` **Status: Pending** (A1/A2/A3/A4/A6/A7 all Verified). A5's own
`Verification plan` defers its implementability half to *"Spike, at
implementation Phase 2"*.

**Clarification requested:** is A5 deliberately a Phase-2 in-implementation
spike (in which case the Prerequisites line needs to except it, or Phase 2 needs
a stated route back if the spike refutes C1's whole-tree clause), or is the
Pending status stale?

**Blocks:** starting at all. The Prerequisites gate is the first checkbox in the
plan, and A5's `If wrong` is *"C1's whole-tree structural negative cannot be
pinned in the house idiom and needs its own form"* — i.e. a Phase-2 refutation
reopens a normative contract mid-implementation. I need to know whether that is
the accepted deal before I cut the branch.

---

## S2 — High

### F3. `0023:C2` (widened: `0023:§source-authority-census`) — the reflective partition oracle is mandated but has no declared in-code home for the assignment

C2 is unambiguous that the mechanism may not be an omit-list: *"the projection
MUST NOT be implemented as a bare omit-list whose complement is 'whatever else
exists'"*, and *"this verb's reflective oracle MUST assert every field of the
resolve success payload is assigned to exactly one group, so an unassigned new
field is a test failure"*.

What C2 does not say is **where the assignment lives** so reflection can read
it. The candidate homes are materially different edits:

- a struct tag on `resolvePayload` (`envelope:"echo"|"plan"`) — reflection reads
  it directly, but it adds a tag to a struct that
  `internal/cli/decision_table_0010_test.go:435` pins at exactly 14 fields and
  whose `json` tags drive a pinned wire key order
  (`decision_table_0010_test.go:422-426`);
- a package-level `map[string]group` keyed by JSON name in the *non-test* code;
- a package-level map in the *test* file, reflecting over `resolvePayload` and
  asserting coverage.

The third is the cheapest and is what "reflective *oracle*" most naturally
reads as — but it puts the partition in a test file while C2 calls it a
normative contract of the verb, and `§source-authority-census` presents the
partition as a doc table whose staleness "the reflective oracle is what keeps
… from going stale". A table in an RDR plus a map in a test is a third copy of
the partition alongside C1's explicit key-set literal.

**Clarification requested:** which of the three is intended, and is the
projection function's own code allowed to *be* the assignment (i.e. reflection
over a projected-struct's field set), or must the assignment be a separate
declaration the projection reads?

**Blocks:** the Phase 1 mechanism choice. A2 defers pointer-`omitempty` vs
projected-struct to implementation; that choice is not independent of this one,
because a projected-struct implementation makes the field set itself the
assignment and a pointer-nilling implementation does not.

### F4. `0023:§implementation-plan` Phase 3 (widened from `0023:D-naming`'s `flowResolveExtendedDesc` anchor) — editing the help body dirties a committed generated file the build gates on, and Phase 3 does not say to regenerate it

Phase 3 says *"`flowResolveExtendedDesc` gains the flag under 'Reading a
successful plan'"*. That string is a `withExtendedHelp` annotation, and
`internal/cli/docs.go` generates `docs/cli-reference.md` from exactly those
extended bodies. `make check` runs `docs-check` (Makefile:34, :105), which
compares regenerated output against the committed copy;
`docs/cli-reference.md` already contains the "Reading a successful plan"
body.

So the Phase 3 edit as written turns `make check` red until
`docs/cli-reference.md` (and `llms.txt`) are regenerated. Phase 3 names
`docs/cli-output-contract.md` — the hand-written one — and not the generated
one.

**Clarification requested:** confirm Phase 3 includes `make docs` and the
regenerated `docs/cli-reference.md` in the same change. Also confirm the flag's
`Short` usage string (`cmd.Flags().Bool("plan-only", false, "<usage>")`) is
free-form — it lands in the generated `Flags:` block too, and no element in this
record fixes its wording, while `0023:D-naming` fixes the flag *name* down to
the rejected spellings.

**Blocks:** the Phase 3 commit passing CI, and the wording of a user-visible
string that ships into a committed artifact.

---

## S3 — Medium

### F5. `0023:C1` — "the carried keys keep their default-mode relative order" is a wire-only property, but the clause sits above the text-mode paragraph without scoping

C1 says *"the carried keys keep their default-mode relative order"* and,
further down, *"The projected text lines MUST be a subset of the default-mode
text lines, byte-identical per line."*

Wire order is struct declaration order (`clierr.WriteJSONLine` →
`json.Encoder`). Text order is **alphabetical** —
`internal/cli/respond/text.go::flatten` calls `sort.Strings(keys)` on the
decoded map. `§desk-trace` states this coexistence explicitly and correctly
("Two orderings coexist and are not the same ordering"). C1 itself does not,
and an implementer reading C1 top-to-bottom could reasonably write a text
oracle asserting *ordered* subsequence rather than *set* subset.

**Clarification requested:** confirm S5's "subset" is set-membership per line
(order-free), not subsequence. `§desk-trace` says yes; C1 does not say either.

**Blocks:** the S5 text-subset oracle's assertion form. Low risk of a wrong
implementation (desk-trace disambiguates), but it is the kind of thing that gets
written wrong at 9am on Monday from the contract alone.

### F6. `0023:C1` — the "STRICTLY SHORTER" width clause is asserted on which encoding, and against what baseline?

C1: *"the projected encoding STRICTLY SHORTER than the default"*. Testing
Strategy S1 makes the same claim and grounds it in *"the five echo keys always
render, `model`/`outcome` being non-`omitempty` and the three containers
rendering `{}`/`[]`"*.

Two ambiguities an implementer hits immediately:
1. **Which encoding** — the `data` payload's own marshalled bytes, or the whole
   NDJSON envelope line (`{"type":"ok","data":{…}}`)? `§testing-strategy` S1's
   fixture cites *"the full NDJSON line measures 290 B and its projection
   152 B"*, implying the envelope; `§performance-expectations` cites
   `708→146 B` for the same class, implying the payload. Both are cited as
   normative-adjacent.
2. **Text mode too, or JSON only?** S1 is scoped `--as=json`; S5 asserts subset,
   not width. If text-mode width is unasserted, say so.

**Clarification requested:** fix the width oracle's unit (payload bytes vs
envelope line) and its mode scope.

**Blocks:** writing S1's width assertion — the one clause the RDR names as the
sole guard against premortem P-12, so it should not be measured on an
ambiguous unit.

### F7. `0023:C2` / `0023:D-wire-byte-format` — the always-keep core is declared but no oracle in the record enforces it

C2 declares an ALWAYS-KEEP core (`rule`, `escaped`, `escape_class`,
`revision`) and calls it *"projection-invariance, not unconditional
presence"*, binding *"every FUTURE projection mode"*.
`§risks-and-mitigations` leans on it: *"C2's always-keep core … binds every
FUTURE projection mode now, so the safety invariant survives any widening of
the axis."*

None of S1–S5 asserts it. S2 pins an explicit literal key set that happens to
contain the core today, but that is a today-assertion: a future enum mode that
drops `escaped` would add its own key-set literal and S2 would stay green on
the boolean mode.

**Clarification requested:** is the always-keep core deliberately unoracled at
this RDR (enforced only when a second mode is proposed), or should Phase 2
carry a sixth oracle asserting the core survives *any* projection the verb
implements — e.g. reflective over the projection function rather than over one
mode's output?

**Blocks:** the Phase 2 oracle count. The RDR says "The Phase 2 oracle battery
is the strategy; done = every oracle green", so I need to know whether the
battery is five or six.

### F8. `0023:C1` — "byte-identical refusal envelopes" ± the flag: is `--plan-only` accepted *at all* on a refusing `flow resolve` run?

The `§disposition` table row *"`--plan-only` on `flow resolve`, refusal | unchanged | refusal, byte-identical (A6) | — | loud, flag-blind"* is clear that the refusal path is flag-blind, and A6 grounds it (`respond.Fail` has no echo member).

What is unstated: cobra parses flags before `RunE`, so `--plan-only` is accepted
on every `flow resolve` invocation including refusing ones, and the projection
code simply never runs. That is the obvious reading and almost certainly right.
But C1's *"MUST NOT change … how the run fails"* combined with MVV step 4
(*"unrecognized outcome, ± the flag → byte-identical refusal envelopes"*)
leaves open a reading where a refusal-path assertion is needed that the flag was
*ignored* rather than *unreachable*.

**Clarification requested:** confirm the refusal-path guarantee is structural
(the projection site is after the last `respond.Fail` return in
`runFlowResolve`, so no refusal can reach it) rather than behavioural, and that
the differential oracle's refusal arm is corroboration in the same sense C1
calls the command-error run corroboration.

**Blocks:** whether Phase 1 needs a guard at the projection site or the
placement alone carries the contract. Cheap either way; I want it written down
so the next reader does not add a defensive branch.

---

## S4 — Low / confirm-and-move-on

### F9. `0023:C1` (widened: `internal/cli/decision_table_0010_test.go:435`) — the pointer conversion must keep `NumField() == 14`; A3 flags it but no contract does

`A3`'s Evidence carries the caveat: *"`decision_table_0010_test.go` pins
`resolvePayload` at exactly 14 struct fields — the A2 mechanism keeps the count
(pointer conversion adds no field)"*. I confirmed the pin
(`decision_table_0010_test.go:434-437`), and the wire key-order pin at
`:422-426` (13 keys, unescaped case).

This is correct for the pointer-`omitempty` arm of A2. It is **not** correct for
the projected-struct arm A2 leaves open and `§risks-and-mitigations` explicitly
contemplates (*"a projected-struct implementation (A2) duplicates
`resolvePayload`"*) — a second struct does not move the count, so that arm is
fine too, but a *shared* struct with an added `Projected bool` or similar would
break it.

**Clarification requested:** none blocking — just confirm the field-count pin is
understood as a hard constraint on Phase 1 rather than a test to be updated.
`0023:A3`'s whole claim is that no predecessor oracle moves; a bumped `14` would
falsify it.

### F10. `0023:D-naming` — "long form only (no shorthand)" is stated twice but no oracle asserts it

C1: *"long form only (no shorthand)"*. `D-naming`: *"boolean, long form
only"*. The `--all` precedent has a shipped oracle for exactly this
(`TestReq36And50_NextRegistersABooleanAllFlagDefaultFalseWithNoShorthand`,
cited at `§references`). No S-scenario here asserts type, default, or absent
shorthand.

**Clarification requested:** should S4 (or a small S6) assert the positive
registration shape on `flow resolve` — boolean, default `false`, `Shorthand ==
""` — mirroring the 0011 precedent? S4 as written asserts only the *set* of
registering commands, not the flag's shape on the one that has it.

**Blocks:** nothing hard; a one-line ask that closes the gap between the two
`D-naming` sentences and the oracle battery.
