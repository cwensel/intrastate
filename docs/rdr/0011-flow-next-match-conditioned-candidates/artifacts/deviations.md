# Deviations — RDR 0011

Gaps found while writing the Phase 1 oracles, each grounded against the RDR
text, `req-list.md`, the 0005 predecessor artifacts, and the source, and each
resolved on the reading the evidence supports. This run is unattended, so
nothing here halted the work: every entry names its Type and Status and the
tests proceed on the recorded reading.

---

## DEV-1 — A13's "one reachable producer" of `uncomparable` is not reachable through a single model

- **Type**: `SPEC-UNDER`
- **Status**: `mechanical translation`
- **Bears on**: REQ-24, REQ-66, REQ-84, REQ-112, `0011:S7`, `0011:A13`

**The clause.** `0011:A13` establishes that of `evaluateAtom`'s four
producers of `ReasonUncomparable`, exactly ONE is reachable from `flow next`:
"the seam-unevaluable arm, `internal/guard/grammar.go::Evaluator.Evaluate`
returning `GuardUnevaluable` for a present value the operator cannot parse
(`lt`/`lte`/`gt`/`gte` on a non-integer, `in`/`contains` on a non-array
literal or held value…)". C3 then requires the fixture and forbids the other
three: "the other three producers are unreachable and MUST NOT be fixtured,
since a fixture that forces one tests `internal/resolve`, which this RDR does
not change" (REQ-66).

**What the source shows.** Every value channel into the assembled view is
validated against the DECLARATION it crosses, and both refusals precede any
probe:

- `--tag size=notanint` over an `int`-declared key → `flow-tag-invalid`
  ("the value for `size` does not conform to its declaration: \"notanint\" is
  not an int"). Verified live.
- `--tag marks=null` / `marks=notarray` over a `set`-declared key →
  `flow-tag-invalid` ("the tag `marks` is set-valued and takes a JSON array
  literal"). Verified live. This closes `parseHeldSet`'s `null` route, whose
  own doc comment says "only a runtime caller can supply it".
- `--write size=notanint` / `--write marks=null` → the same refusals under
  `flow-write-invalid`. Verified live.

And the loader closes the other side: `internal/table/normalize.go:94`
applies RDR 0003's operator/kind matrix to GUARD blocks, so
`[rule.guard.all.size] gt = "3"` is refused `malformed_predicate_atom`
("operator gt does not accept kind scalar") on any kind but `int`.

So within ONE model the two halves are mutually exclusive: a key an operator
can be authored against is a key whose channel refuses the value that would
make it unparseable.

**The reading taken.** The artifact is CALLER-OWNED and is not a per-model
store — `flowbind`'s reader loads whatever the bound path holds, and RDR
0005's own harness rule is that no test may assume the artifact's on-disk
format. So the value is ESTABLISHED under one model where it conforms and
READ under another where the same key carries a stricter declaration:

- `flowLooseWriterModel` declares `size` a `scalar`; `--write size=notanint`
  conforms and lands through the production write path.
- `flowUncomparableGuardModel` declares the same `size` an `int` and guards
  `gt = "3"`. The reader establishes `notanint`, so the key is PRESENT at a
  value `strconv.Atoi` cannot parse and the seam answers `GuardUnevaluable`.

Verified live: the two-model route produces the disposition; a single-model
route cannot. Both hops go through the CLI, so the 0005 harness rule holds.

**Why this is mechanical rather than an author decision.** A13 names the
producer by its SEAM behaviour ("a present value the operator cannot parse"),
not by a construction; the reading supplies the only CLI-reachable
construction that reaches it. It reaches the producer A13 names and does not
touch the three A13 calls unreachable — no `TagSet` is constructed, no nil
seam is forced, no non-boolean existence literal is authored. `internal/resolve`
is untouched, which is the constraint C3 attaches to the prohibition.

**What a later phase could overrule.** If a successor decides a cross-model
artifact read is itself out of contract, the alternative is that
`uncomparable` has NO CLI-reachable fixture on this build and S7's clause is
unattestable until some model-scoped store lands — in which case A13's
"If wrong" branch applies verbatim: "`uncomparable` never appears on the
CLI's list and the `{key, reason}` shape carries one live reason fewer; the
shape still stands on `absent` vs `not-evaluated`."

---

## DEV-2 — REQ-100's "two match tags on the dead row" is not observable from the CLI payload

- **Type**: `SPEC-UNDER`
- **Status**: `mechanical translation`
- **Bears on**: REQ-8, REQ-9, REQ-80, REQ-99, REQ-100, `0011:S3`, `0011:A12`

**The clause.** `0011:S3` states "the oracle asserts two match tags on the
DEAD row only — the pairing's other expanded row collapses to one tag and is
live, so asserting two there would assert something false."

**What the payload carries.** A `flow next` candidate carries `rule`,
`outcome`, `required`, the undecided list, `next`, `writes`, `clear`, and
(under `--evaluate-gates`) `gates`. It does NOT carry the row's match
pattern — Failure Modes says so outright: "the payload carries no `match`
block in either mode". Both expansions of an `eq`+`in` pairing also share ONE
rule id, so the two rows are indistinguishable by identity as well.
Confirmed live: `flowMatchClassesModel` reports `dead-row` twice, byte-identical.

Counting tags would therefore mean reading `Row.KernelRow().Match` in the
test — reaching past the CLI surface into the normalizer, which no oracle in
this suite does and which the 0005 harness header explicitly rules out.

**The reading taken.** The tag-count asymmetry is asserted through its
observable CONSEQUENCE — the DISPOSITION, which is what C1 actually
specifies:

- `phase` absent → `dead-row` reported TWICE (both expansions survive; the
  filter omits the key's tags together, never splitting a key — REQ-7).
- `phase=alpha`, the value where the `in` member COINCIDES with the `eq`
  literal → `dead-row` reported ONCE. The collapsed one-tag expansion is
  live; the two-tag expansion is dead and the kernel's conjunction excludes
  it. This is REQ-100's asymmetry, and asserting two survivors here would
  assert exactly the false thing S3 warns about.
- `phase=zeta`, a value on NEITHER side of the pairing → `dead-row` absent
  entirely; both expansions fail.

The three arms together pin the pairing's shape as tightly as the tag count
would, and they pin it against the KERNEL's conjunction rather than against a
CLI-side count — which is the property REQ-10 forbids the CLI from having.

**Fixture note.** The dead row is authored A12's way (`eq = "alpha"`,
`in = ["alpha", "zeta"]`), NOT with the `eq` literal outside the set. The
outside-the-set form expands to two rows that are BOTH dead, so it has no
live control and cannot exhibit the collapse REQ-100 turns on.

---

## DEV-3 — REQ-109's key-set comparison cannot call the kernel's `assemble`

- **Type**: `DEPENDENCY-LIMIT`
- **Status**: `mechanical translation`
- **Bears on**: REQ-109, `0011:S5`, `0011:A11`

**The clause.** `0011:S5`: "an oracle MUST compare `assembledView`'s key set
against the kernel's `internal/resolve/resolve.go::assemble` over the same
inputs, modulo the kernel's `recognized` key".

**The limit.** `internal/resolve/resolve.go::assemble` is UNEXPORTED, and so
is `TagSet`'s key set — `TagSet` exports only `Lookup` and `Len`, neither of
which enumerates. Calling it from `internal/cli` would require exporting
something under `internal/resolve`, and S5 forbids that in the same
paragraph: "the kernel is untouched — the ONE mechanically checkable form of
that claim is MVV 9's `git diff --stat internal/resolve` empty". The two
halves of S5 cannot both be satisfied literally.

**The reading taken.** The comparison is made against `assemble`'s
SPECIFICATION rather than against the function — owned ∪ observed ∪
`recognized`, which is exactly what A11 says it builds from the `Input` the
CLI hands it — reconstructed from the payload's own `owned` and `observed`
maps. The oracle asserts:

1. `assembledView` does NOT carry `recognized` (the one key the two builders
   differ on, and unnameable by any atom since `normalize.go` lifts the
   match-block `recognized` atom into `Row.Outcome`);
2. a key `--tag` supplied IS in the CLI's view; and
3. the agreement in the direction that matters — a present-and-equal key
   DECIDES and a present-and-unequal one EXCLUDES. A key the kernel's view
   LACKED would fold into `no_match` and drop the row silently, so neither
   observation is producible under a key-set disagreement.

**What is preserved.** S5's stated purpose survives: "a refactor that stops
threading one `owned` slice to both consumers fails here rather than
silently." (3) is what breaks under such a refactor — the CLI would keep a
match atom the kernel's view lacks, and the row would vanish.

**What is lost.** The oracle is behavioural rather than set-theoretic, so it
discriminates a disagreement on a key some row MATCHES on, not on an
arbitrary key. That narrowing is harmless: A11's claim is scoped to "every
key a row's match atom can name" in its own statement.

---

## DEV-4 — the S8 `flow resolve` plan arm needs a second fixture

- **Type**: `TEST-FIXTURE`
- **Status**: `mechanical translation`
- **Bears on**: REQ-19, REQ-20, REQ-93, REQ-116, `0011:S8`

**The clause.** `0011:S8` asks one model to carry both halves: the three
`flow resolve` arms (plan / role-unbound / refusing, REQ-93) AND "the
BREAKING arm: add a second ordinary row for the same outcome that does not
match on the key" (REQ-116).

**The conflict.** With the reader bound and answering `mode=fast`, BOTH `go`
rows match — the match-only row on `mode`, and the breaking arm's `plain-row`
on `status`. The kernel then refuses `ambiguous_match` rather than planning,
so the plan arm is unreachable on the two-row model. This is a property of
the kernel's exactly-one selection (`0005:D-selection-predicate`), not of
anything this RDR changes.

**The reading taken.** Two fixtures, one per question S8 asks:

- `flowMatchOnlyOwnedSoloModel` — one `go` row. Asks what the term BUYS: the
  plan arm, and the two refusal arms that name the remedy the old
  `flow-no-match` hid.
- `flowMatchOnlyOwnedModel` — two `go` rows. Asks what the term COSTS: the
  breaking arm, where the row `resolve` selects does not itself need the key
  and the union-over-rows demand pulls the reader in anyway.

Both are S8's; splitting them is the only way to observe both, and neither
weakens the other. The refusing-reader arm rides the solo model too
(`flowRefusingSideReaderModel`), for the same reason.

**A15's own limit is respected.** The exit-3 oracle asserts the EXIT CODE and
the not-a-`flow-no-match` property, never a single code — A15(c): "another
refusal mechanism would take a different exit-3 code, so what is verified is
the exit code and the not-a-`flow-no-match` property, not a single code."

---

## DEV-5 — REQ-58 and REQ-59 pass vacuously before the re-homing lands

- **Type**: `SPEC-UNDER`
- **Status**: `mechanical translation`
- **Bears on**: REQ-58, REQ-59, `0011:S6`, `0011:C3`

**Found by the red gate.** Both oracles scan the three census files for a
defect in the RE-HOMED reads: a discarded comma-ok (REQ-58) and a repointed
`stringsAt` (REQ-59). Before the rename those reads do not exist, so the scan
finds nothing and passes — for exactly the reason C3 warns against: "a green
suite MUST NOT be cited as evidence that this contract was implemented."

**The reading taken.** Both now open with a VACUITY GUARD that `t.Fatal`s
while any census file still reads `"unresolved"`, on the ground that the
clause is about the re-homed reads and is not assertable until they exist.
Both are consequently RED today and counted under `red_confirmed`.

This is the same idiom the repo already ships for absence oracles
(`TestReq69_NoPlanFlagShipsOnAnyVerb` guards on the four verbs existing;
`TestReq2_FlowGroupDoesNotAbsorbLintDumpOrParse` on the group existing), so
the resolution is the house pattern rather than a new one.

---

## DEV-6 — REQ-117's reader-set agreement is unobservable on the two-row fixture

- **Type**: `TEST-FIXTURE`
- **Status**: `mechanical translation`
- **Bears on**: REQ-17, REQ-117, `0011:S8`, DEV-4
- **Found in**: Phase 2, implementing the demand-set term

**The clause.** `0011:S8`: "Both verbs' `readers` over the same model and
outcome are EQUAL (one `invokedReaders`, one view), asserted on a fixture
whose reader serves several keys."

**What the implemented build shows.** Phase 1 wrote
`TestReq17And117_BothVerbsAgreeOnTheReaderSetOverOneModelAndOutcome` over
`flowMatchOnlyOwnedModel` — the TWO-row fixture — with `mode=fast` and
`--outcome go`. Once the demand-set term lands and `read.side` is invoked,
`mode=fast` makes BOTH `go` rows match: `mode-row` on `mode` and
`plain-row` on `status`. The kernel then refuses `ambiguous_match`
(observed live: `flow-ambiguous-match`, naming both rules) and `flow
resolve` emits a REFUSAL envelope, which carries no `readers` list at all.
The oracle's comparison target does not exist on that input.

This is the SAME collision DEV-4 already recorded in the other direction —
"With the reader bound and answering `mode=fast`, BOTH `go` rows match …
the kernel then refuses `ambiguous_match` rather than planning, so the plan
arm is unreachable on the two-row model" — applied to the `readers`
comparison rather than to the plan arm. DEV-4 split S8's arms across two
fixtures for exactly this reason and simply did not carry the split through
to this one oracle.

**The reading taken.** The oracle moves to `flowMatchOnlyOwnedSoloModel`,
which satisfies S8's stated fixture requirement verbatim: its `read.side`
declares `keys = ["mode", "extra"]`, so the agreement IS "asserted on a
fixture whose reader serves several keys" — the property S8 names — and
`--outcome go` has exactly one row, so `flow resolve` plans and reports a
`readers` set to compare. Verified: both verbs report `[side state]`.

Nothing is weakened. The assertion still compares `next`'s `readers`
against `resolve`'s over one model and outcome, still fails a `next`-only
term, and still asserts `side` is present on the `resolve` side — the two
halves REQ-17 and REQ-117 name. The two-row model keeps its own oracle,
`TestReq20And116_TheBreakingArmIsPinnedAsAKnownCost`, which is the arm S8
authors it for.

**Why mechanical.** The clause names the fixture PROPERTY ("whose reader
serves several keys"), not the fixture identity, and the solo model has
that property. The change is a fixture repoint within the corpus DEV-4
already established, with no contract wording touched.

---

## DEV-7 — REQ-125's element-type check is vacuous over an empty `unknown`

- **Type**: `TEST-FIXTURE`
- **Status**: `mechanical translation`
- **Bears on**: REQ-125, `0011:G-cross-cutting`
- **Found in**: Phase 2, against the implemented payload

**The clause.** REQ-125: "Version marker: none minted — the payload field
rename `unresolved` → `unknown` IS the break". Phase 1's
`TestReq125_NoVersionMarkerIsMintedAndTheOldFieldIsGone` adds a third
assertion beyond the clause's two: that `unknown` decodes as a list of
OBJECTS, not of strings, spelled `if _, ok := stringsAt(c, "unknown"); ok`.

**What the source shows.** `flow_harness_0005_test.go::stringsAt` returns
`(out, true)` whenever the value is a `[]any` every element of which is a
string — and an EMPTY `[]any` satisfies that vacuously, since the loop
never runs. An empty JSON array carries no element type at all.

The oracle runs over `flowAbsentMatchKeyModel`, whose `recognized-only`
row matches only on `recognized` and writes `status`, which the reader
establishes. Its `unknown` is therefore `[]` — verified on the wire:

```
{"rule":"recognized-only","outcome":"hold","required":["status"],
 "unknown":[], ...}
```

So `stringsAt` returns `([], true)` and the assertion fires against a build
whose payload is exactly what the contract requires. It would fire the same
way against a `[]string` build, a `{key,reason}` build, and a build with no
field at all — it discriminates nothing on that candidate.

**The reading taken.** The check is scoped to a candidate that actually
carries an entry, by `continue`ing past an empty list. `absent-row` on the
same fixture carries `{wanted, absent}`, so the assertion still runs on
real data every time and still fails a `[]string` build.

The two assertions REQ-125 itself states — no `schema`/`version`/
`payload_version`/`$schema` field, and no surviving `unresolved` field
alongside `unknown` — are untouched and green.

**Why mechanical.** The clause's own obligations are unchanged; only the
supplementary type check is scoped away from an input that cannot express
it. This is the same shape as DEV-5's vacuity guard, applied to the
opposite polarity: DEV-5 added a guard so an assertion could not PASS
vacuously, and this one adds a guard so an assertion cannot FAIL
vacuously.

---

## Non-deviations, recorded so a later phase does not re-open them

- **The C3 census reconciles exactly.** REQ-54's five `file:line` citations
  were checked against the tree: `flow_next_0005_test.go:103,163,404`,
  `flow_adversarial_0005_test.go:446`, `flow_mvv_0005_test.go:66`, plus
  `:166` as the `%#v` argument C3 names as not-a-read. No drift. Both REQ-54
  and REQ-60 are nonetheless written against the SYMBOL and the COUNT rather
  than the line, per `§source-location-discipline`.

- **A4's MOVES-UNDER-`--all` = 0 holds.** No 0005 `next` oracle needed
  re-homing under `--all` while writing these tests, which is what A4
  predicts. `TestReq53And111` asserts the negative directly, so a move a
  later phase finds necessary REFUTES A4 and must be recorded here rather
  than absorbed.

- **No round-trip / inverse invariant is asserted.** The mini-check cue
  register says "Cues absent: round-trip / fidelity — this RDR defines no
  import/export, parse/deparse, or serialize inverse." The MVV's invariant is
  the candidate set's exact membership in both directions instead, which is
  what `TestReq87And95` asserts (every row the state can take is reported,
  every row it excludes is absent) rather than an exit code.

- **A-12 (Q1) and A-13 (Q2) are taken as `req-list.md` records them.** The
  `--all` filter is per-ATOM on `Block == BlockMatch`, so a GUARD atom on the
  same absent key still contributes; and `--all` does not re-widen the demand
  set back to 0005's. Both are asserted
  (`TestReq7And103` and `TestReq18And21And115` respectively), so a later
  phase that overrules either will fail a named oracle rather than drift.
