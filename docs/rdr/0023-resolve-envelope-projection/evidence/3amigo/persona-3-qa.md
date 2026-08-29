Model: claude-opus-5[1m]

# Persona 3 — QA / Tester

Owned starting set read by id: `0023:C1`, `0023:C2`, `0023:S1`–`0023:S5`,
`0023:MVV`. **Widened** to `0023:§oracle-discriminability`,
`0023:§phase-2-oracles`, `0023:§desk-trace`,
`0023:§performance-expectations`, `0023:§failure-modes`, and to
`0023:A2`/`0023:A3`/`0023:A5` plus JDR 0002 §D1. What sent me: S1's
"identical invoked-reader set" clause names an observable that C1 projects
away, and no owned element says how it is read; S3's "reflective" oracle
names no assignment source inside C2; and S4's totality clause depends on
command-tree construction facts stated nowhere in the owned set. A silence
has no line range.

Source checks run against `/Users/cwensel/sandbox/newcoinc/intrastate`
(HEAD, `go test`/`go build` probes on the real tree) — cited inline.

---

## HIGH

### H1 — `0023:S1` / `0023:C1`: the "identical invoked-reader set" clause names no observation channel that survives the projection

C1's report-only battery and S1's Expected both require an oracle asserting
"an identical invoked-reader set" over the same request ± `--plan-only`. The
only shipped channel for that set is the payload's `readers` field
(`internal/cli/flow_resolve.go:228`, `Readers: readerIDs(readers)`; the
existing precedent oracle reads it exactly that way —
`internal/cli/flow_fixtures_0011_test.go:1769 readersOf`, which does
`stringsAt(data, "readers")` and `t.Fatalf`s when absent). But `readers` is
an ECHO field under `0023:C2` and is ABSENT under the flag by C1's own
key-set clause and by `0023:MVV` step 2. So on the projected side of the
differential there is nothing to compare: the one available reader oracle
`Fatal`s on the projected run rather than passing.

No element names a substitute channel. Readers are real accessor
invocations against artifact roles (`internal/cli/flow_exec.go:232
runReaders`), so a side-effect channel (touch-witness artifacts, an
injected binding spy) is conceivable, but nothing in C1, S1, MVV, or
`§phase-2-oracles` picks one, and A5's spike covers the *registration* walk,
not this.

**Test prevented:** the differential reader-set assertion of S1 — i.e. the
one assertion that would catch "a later change skips work whose only
consumer is a projected-away field", which C1 explicitly declares a breach.
Without a named channel the implementer will either silently drop the
clause (leaving C1's most load-bearing report-only guarantee unenforced,
contradicting "report-only is ORACLE-ENFORCED, not aspirational") or invent
a channel the RDR never sanctioned.

**Blocks:** the report-only clause of `0023:C1` and Phase 2's differential
oracle.

---

### H2 — `0023:S4` / `0023:C1`: the totality clause is unsatisfiable as written against a test-constructed tree, and S4's own negative control is blind to the gap

C1 says the walk "descends every child of the root with no name-based skip
and no `Hidden` gate, INCLUDING the auto-generated help and completion
commands", and S4's Expected repeats "`help` and `completion` included".

Measured on HEAD: a tree built by `NewRootCmd()` contains **no `completion`
command at all**. Probe (temp test, since removed):

```
child: "docs" hidden=true
child: "flow" hidden=false
child: "help" hidden=false
child: "lint" hidden=false
child: "version" hidden=false
```

`help` is present only because `internal/cli/help_all.go:135` force-calls
`root.InitDefaultHelpCmd()`. `completion` is created lazily by cobra inside
`Execute()` via `InitDefaultCompletionCmd()`, which nothing in
`internal/cli/root.go` calls at construction. Adding an explicit
`root.InitDefaultCompletionCmd()` makes the walk reach
`completion bash|fish|powershell|zsh` (verified by probe), so the oracle
*is* writable — but only via a force-init step no element states.

Worse, S4's stated negative control is "a control asserting the walk
actually reaches `help` guards against a walker that silently narrows."
`help` is present by default; `completion` is not. The control therefore
passes on precisely the tree that omits the branch C1 calls load-bearing.
The oracle can be green while asserting a *weaker* negative than C1
requires — the exact defect C1 spends a paragraph rejecting
`help_all.go::walkCommandTree` for.

**Test prevented:** the total-tree registration walk with a control that
actually discriminates. As specced, S4 ships as a walk over four commands
plus `help`, indistinguishable from the enumerated sweep C1 says it must be
"deliberately stronger" than.

**Blocks:** C1's non-registration clause and the WHOLE-means-total
paragraph.

---

## MEDIUM

### M1 — `0023:S3` / `0023:C2`: "reflective partition-completeness" names no artifact to reflect *against*, so the oracle can be written tautologically

C2 requires an oracle asserting "every field of the resolve success payload
is assigned to exactly one group" and forbids "a bare omit-list whose
complement is 'whatever else exists'". But C2 states the assignments only in
prose ("the ECHO group is the model reference, the observed tags, the
assembled owned view and the invoked reader identities, and the requested
outcome"). JDR 0002 §D1 likewise says "enforced by structure, not prose"
without naming the structure. Neither C2, S3, nor
`§oracle-discriminability` mandates a declared in-code assignment table
(e.g. two exported/`var` field-name lists) that the reflection compares
against.

If the implementer derives the ECHO set from the projection function itself
— which is the path of least resistance under A2's pointer-nilling
mechanism, where "which fields are echo" is expressible as "which fields are
pointers" — S3 reduces to "the code agrees with itself" and still passes
when a new field is added as a pointer, or when a new non-pointer field
silently joins the projected width. S3's negative control ("add a field to
the payload without a C2 side; S2 and S3 must both go red") does not
discriminate this: under a self-derived assignment, a plain new field
*would* be red on S2 (authored literal) but green on S3, so the control's
"both" condition is unverifiable at authoring time.

**Test prevented:** an S3 that fails on an *unassigned* field independently
of S2. Today S3 is at risk of being S2 with extra reflection.

**Blocks:** C2's enforcement clause and the anti-omit-list mandate.

---

### M2 — `0023:S1` / `0023:C1`: the width clause has a stated pass criterion for JSON only; the same "whole reason for existing" is unenforced in text mode

C1 asserts "the projected encoding STRICTLY SHORTER than the default" and
calls it "the flag's whole reason for existing … itself oracle-enforced".
S1 scopes that assertion to `--as=json`. S5, the only text-mode scenario,
asserts subset-and-byte-identity per line but never a strictness/width
criterion, and C1's text paragraph likewise says only "subset".

A subset is satisfied by the *equal* set. So the P-12 hazard C1's width
clause exists to catch — a later plan-group field growing the payload with
key-set and partition oracles green — has an enforced floor in JSON and no
floor in text. The `§oracle-discriminability` row for S5 offers the control
"drop a plan-group line under the flag; S5 must go red", which tests the
opposite direction.

Note the criterion *is* writable: measured on HEAD the default text
rendering of the normative fixture is 15 lines, projected 9 (verified by
running the built binary), and `respond/text.go::flatten` emits a line for
every leaf including empty containers, so the six echo lines always render.
The RDR just never states the threshold.

**Test prevented:** a text-mode width regression oracle. As specced, a
text-mode projection that saved nothing passes S5.

**Blocks:** the width half of C1's report-only clause in text mode.

---

### M3 — `0023:§phase-2-oracles` vs `0023:S1`–`0023:S5`: the completion criterion counts six oracles; the Testing Strategy enumerates five

`§testing-strategy` opens "done = every oracle green and the MVV run
recorded", and `§phase-2-oracles` lists six: whole-tree registration,
differential report-only, reflective partition-completeness, explicit
projected-key-set, text-subset, **and absent-not-null**. The scenario
elements are S1–S5; there is no S element for absent-not-null (it is folded
into S2's Expected as a trailing "omitted keys absent — never null or empty
placeholders"), and it has no row in `§oracle-discriminability` and no
negative control.

Because the phase's done-criterion is "every oracle", the battery's cardinal
count is ambiguous at implementation time: an implementer who works from the
S list ships five and can call Phase 2 complete; one who works from Phase 2
ships six.

**Test prevented:** none directly, but the ambiguity means the
absent-not-null property ships with no discriminability analysis — and it is
the one property A2's spike identifies a live hazard for (bare non-pointer
`omitempty` drops `"owned":{}`, reproduced in
`evidence/spikes/a2-encoder-mechanism.md`). A folded assertion with no
negative control is exactly where that hazard re-enters.

**Blocks:** the Phase 2 done-criterion.

---

## LOW

### L1 — `0023:S1`/`0023:S2` cite a "normative fixture" whose byte counts do not appear in the cited artifact

S1 states "the full NDJSON line measures 290 B and its projection 152 B" and
pins both to `evidence/spikes/a2-encoder-mechanism.md`. The artifact carries
the two literal lines and they do measure 290/152 (verified), but neither
number nor the word "normative" appears in the file — `rg '290|152|normative'`
over it returns nothing. A byte-count assertion written from the RDR is
therefore checkable only by re-deriving `len()` from the spike's prose
blocks; the traceability the "normative fixture" language promises is
one hop shorter than stated. (A1's table at
`evidence/spikes/a1-byte-width.md:25` does carry `290 | 152` for shape S1,
which is where an implementer would have to look.)

**Test prevented:** none — the fixture bytes are recoverable. Flagged
because two scenario Expecteds route their exactness claim through a
citation that does not contain the cited figures.

### L2 — `0023:§performance-expectations`: "Wall time unchanged" carries no threshold and no oracle

The clause asserts the 5ms call is unaffected, but names no measurement, no
tolerance, and no scenario. No S element covers it and `§phase-2-oracles`
lists no timing oracle.

**Test prevented:** any performance-regression assertion. Given C1's
"projection is key deletion … no second marshal path", this is plausibly
intentional (the claim is structural, not empirical) — but as written it
reads as an expectation a tester is expected to check and cannot.

### L3 — `0023:MVV` step 5 exercises one sibling verb behaviourally where `0023:F1` names three

C1's non-registration clause covers "any other command"; F1 names `next`,
`read-state`, and `set-state` as the visible `command-error` cases. MVV
step 5 runs only `flow next --plan-only`. S4's structural walk is the real
guarantee and C1 says "the behavioural command-error run is corroboration,
never the assertion" — so this is consistent by design. Recording it only so
a reader of F1 does not expect three MVV runs and find one.

---

## Checked and found writable (no finding)

Listed so the absence of a finding is legible as a check, not a gap:

- **S2 key set** — `resolvePayload` (`internal/cli/flow_resolve.go:31-55`)
  carries all nine plan keys; `escape_class` is the only `omitempty` member
  and C1/S2 correctly bind its presence to the same run's default output
  rather than to `escaped`.
- **S1 width clause unconditionality** — verified: `observedTagMap`/`tagMap`
  (`internal/cli/flow_exec.go:739,747`) and `readerIDs`
  (`internal/cli/flow_next.go:582`) all allocate unconditionally, so
  `observed`/`owned`/`readers` always render `{}`/`[]` and the five echo keys
  are always present by default. The clause cannot be vacuous in JSON.
- **S5 determinism** — `respond/text.go::flatten` sorts each object's keys
  (`sort.Strings`, text.go:75), so no map-iteration path exists; the desk
  trace's 15→9 line arithmetic reproduces exactly on the built binary.
- **MVV step 5 exit mapping** — `command-error` / exit 2 for an unknown flag
  is shipped and already asserted elsewhere (`internal/cli/root.go:267`,
  `internal/cli/flow_group_0005_test.go:39`).
- **A3's non-disturbance claim** — `rg plan-only internal/` returns no
  matches on HEAD, so no predecessor oracle moves.
