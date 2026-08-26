Model: claude-opus-5

# 3amigo Persona 3 — QA / Tester (RDR cli/0011)

Scope note: my owned set is S1–S6, MVV, C1–C3. Three findings required WIDENING
past it — into `0011:A7`, `0011:D-undecided-vocabulary`, and the shipped
`internal/resolve` suite — because the untestable part is a silence: the record
names an expected outcome without naming the surface that carries it, and a
silence has no line range inside the owned elements. Each widening says what
sent me.

Baseline established live: `go test ./internal/resolve -run 'TestAdv0007_3'` is
green on `main` at 36e411e, so every claim below about a shipped oracle is
against a passing baseline, not a broken one.

---

## SEV-1 — S5 asserts a shipped oracle "keeps passing" that C1 inverts; two of its siblings break unnamed

**Anchors:** `0011:S5`, `0011:C1` (escape-uniformity paragraph), widened to
`0011:A7` (Status: **Pending**) and `internal/resolve/match_conflicted_test.go`.

S5's Expected says `TestAdv0007_3_EscapeNotNamingTheConflictedKeyStillRescues`
"keeps passing, and gains the sibling case it does not cover." I traced that
fixture and it does not keep passing under C1.

The fixture (`internal/resolve/match_conflicted_test.go:29-72`) builds the
observed view with `gate` supplied twice with differing values, so `gate` is
CONFLICTED. The ordinary row `A` (`gateMatchRow`, line 45) has
`Match: [status=Draft, gate=closed]` and NO guard atoms. Under today's
two-valued `TagSet.matches` (`internal/resolve/resolve.go:174-181`) the
conflicted leg returns false, row `A` drops, `len(selected)==0`, and
`escapeOrRefuse` rescues with `ESCAPE-CLEAN` — which is exactly what the test
asserts (`got.Plan.RuleID == "ESCAPE-CLEAN" && got.Plan.Escaped`).

Under C1, `gate` conflicted ⇒ that atom is `indeterminate`; `status=Draft` is
present and equal ⇒ `match`; row verdict = `indeterminate`; C1 says
"`indeterminate` MUST NOT exclude." Row `A` is therefore a candidate, its guard
is vacuously true, `len(selected)==1`, and `Resolve` returns
`Plan{RuleID:"A", Escaped:false}`. The rescue phase is never entered. The test
fails on both its assertions.

Two siblings in the same file break the same way, and S5 names neither:

- `TestAdv0007_3_DuplicateKeyMakesRowSelectionAFunctionOfSlicePosition` (line
  ~110) asserts `got.Refused()` and `Kind == KindNoMatch` on the same table.
  Under C1 the disposition is a Plan, so `mustResolve`'s refusal assertion
  fails at `t.Fatalf`.
- `TestAdv0007_3_ConflictedKeyIsNotEscapableByNamingIt` (line ~150) is a
  DIRECT INVERSION of C1: its stated contract is "an escape row that NAMES the
  conflicted key must fail the match too... If it did not, an operator could
  rescue the positional plan through the kernel's most permissive path." C1
  says the opposite — an `indeterminate` escape edge "MUST NOT be silently
  disqualified from rescuing." The record retires exactly one oracle by name
  (`TestReq78_MatchPatternStillFoldsAbsenceIntoNonMatch`,
  `internal/resolve/guard_atoms_test.go:930`) and this is not it.

S5's only quantified pass/fail criterion is "the frozen `internal/resolve` suite
passes, with `TestReq78` retired." That criterion is unsatisfiable as written: a
correct C1 build produces a RED suite on three further oracles, and the record
gives no rule for telling "C1 correctly inverted a stale oracle" from "C1 broke
the kernel." The widening to A7 is what makes this blocking rather than
cosmetic — A7 is the assumption that would have enumerated these, and its
Status is **Pending** with Evidence reading "To verify at Stage 6." That
enumeration does not exist in `evidence/`; only forward references to it do.

**Blocks:** the entire S5 regression oracle. I cannot write "run
`./internal/resolve` and require green," because green is the WRONG answer under
C1, and I cannot write the correct expected-failure list because the record
never produced it. Also blocks the S5 sub-oracle "`TestAdv0007_3_Escape...`
keeps passing" — that assertion, written literally, would fail a correct build.

---

## SEV-1 — C1 requires `uncomparable` match facts in `unknown`, but no probe return channel carries them

**Anchors:** `0011:C1` (the `unknown` clause), `0011:S3` (fourth row),
`0011:MVV` step 6, `0011:A1`, `0011:A2`. Widened to `internal/resolve/resolve.go`
`Result`/`Plan`/`Refusal` and `internal/cli/flow_next.go::excluded` — the
record's Existing Infrastructure Audit points at these as the carriers, so I
checked whether they can carry it.

C1: "Every `indeterminate` match atom ... MUST appear in that candidate's
`unknown` list as a `{key, reason}` pair." The Audit splits the sourcing: the
CLI's existing `summarize` atom walk "still fills the absent-key half," and the
kernel supplies the rest. A2 confirms the absent half is already live
(`summarize` at `internal/cli/flow_next.go:184-190` walks `row.Atoms` with no
`Block` test). A3 establishes the `uncomparable` half CANNOT come from the CLI:
`assembledView` (`internal/cli/flow_next.go:113-121`) is a flat `map[string]string`
that collapses a duplicate key to last-wins, so a conflicted key is
indistinguishable from a settled one.

So `uncomparable` must arrive from the kernel probe. It cannot, on the shapes
the record commits to:

- The probe is `resolve.Resolve` returning `Result{Plan *Plan, Refusal *Refusal}`
  (`internal/resolve/resolve.go:390-400`).
- A row whose only defect is an `indeterminate` match atom is, by C1, NOT
  excluded. On a one-row probe with no guard problem, `len(selected)==1` and the
  probe returns a **Plan**.
- `Plan` (`internal/resolve/resolve.go:335-352`) carries RuleID, SourceLocator,
  NextTags, Writes, Revision, Escaped. There is no undecided-facts field.
- `Refusal.Undecided []UndecidedRow` exists (`internal/resolve/resolve.go:372-378`)
  but is documented and populated only "on `guard_unevaluable`", by `gate`
  (`internal/resolve/resolve.go`, the `KindGuardUnevaluable` branch). It is
  unreachable on the success path.

A1 confirms this reading rather than resolving it: it enumerates the probe's
possible dispositions as "a plan, `no_match`, `owned_state_unavailable`, or
`guard_unevaluable`" — and the `uncomparable` case lands squarely in "a plan."
Phase 2 states the intent ("carry the kernel's undecided match atoms into the
candidate's `unknown`") but names no signature, field, or return value; the
Audit row for `TagSet.matches` says "returns a three-valued verdict plus
undecided atoms," which is an unexported helper's shape, not a shape reachable
from `excluded()` through the exported `Resolve` boundary the probe uses.

**Blocks:** MVV step 6's fourth class and S3's fourth row — "supplied twice with
differing values ⇒ candidate with `{key, uncomparable}` in `unknown`." This is
the case the mini-check `oracle` table singles out as "the case the pre-Stage-4
design got wrong" and that S3 says "MUST be asserted." I cannot write it: I have
no named surface to assert the reason token arrives on, so any test I write
either pins an invented API or degrades to asserting only that the row is a
candidate — which does not discriminate `uncomparable` from `absent` and so
fails C3's own discriminability requirement.

---

## SEV-2 — C1 requires a `{key, reason}` pair for gate ids and owned keys, but the closed reason set has no token for either

**Anchors:** `0011:C1` (`unknown` clause), `0011:D-undecided-vocabulary`,
`0011:S2`, `0011:MVV` step 7. Widened to `0011:D-undecided-vocabulary` (a D
element, outside my owned set) because C1 defers the vocabulary there, and to
`internal/resolve/guard.go:67-82`.

C1 puts four distinct fact classes into one `unknown` list, each as a
`{key, reason}` pair: (1) indeterminate match atoms, (2) undecidable guard
facts, (3) "every owned key no invoked reader established", and (4) "absent
`--evaluate-gates`, every gate id."

`D-undecided-vocabulary` closes the reason set: "reuses `0007:C8`'s closed
reason set (`absent`, `uncomparable`) rather than minting one," and C1 repeats
it — "a reason drawn from the closed set `0007:C8` already owns." Source
confirms exactly two members: `ReasonAbsent = "absent"`,
`ReasonUncomparable = "uncomparable"`, with `Reasons()` returning that pair
(`internal/resolve/guard.go:70-82`).

Class (3) fits `absent` cleanly. Class (4) does not fit either. An un-run gate
id is not `absent` — the gate is declared on the row and its id is right there
in the payload; and it is not `uncomparable` — `uncomparable` is defined as "the
key was present and its value was not compared to a verdict," and a gate id is
not a view key with a value at all. It was simply never consulted, because
`--evaluate-gates` was not passed. That is a third sense, and the vocabulary is
closed against minting one.

Nothing in the record picks. The mini-check `disposition` table (`§mini-checks`)
enumerates ten input classes and assigns a reason token to the match classes and
to "owned key no reader established" (`{key, absent}`) — the gate rows say only
"gates run: yes/no" and never assign a token. S2 asserts `unknown` is present as
`[]` on a fully-resolved candidate but never asserts a gate id's shape. Today the
gate id goes in as a bare string with no reason at all
(`internal/cli/flow_next.go:200-207`).

This is not academic: the four `unresolved`-key reads C3 sends for mechanical
re-homing include the gate-id assertions. C3 says so explicitly — "four, at the
presence check and the gate-id assertions." I confirmed the four textual reads at
`internal/cli/flow_next_0005_test.go:103, 163, 166, 404`; the 163/166 pair is
`TestReq41And46And57_WithoutTheFlagGateIdsAreUnresolvedFactsNotResults`, which
asserts `slices.Contains(unresolved, "approval")` on a bare string list. Re-homing
it to `{key, reason}` requires knowing the reason, and C3 calls the re-homing
"mechanical."

**Blocks:** the re-homing of
`TestReq41And46And57_WithoutTheFlagGateIdsAreUnresolvedFactsNotResults`
(2 of C3's 4 named reads) and the `flow_mvv_0005_test.go:66` read. I cannot
write the post-change assertion without a reason token to expect, and C3's
"mechanical" framing means no one else is expected to decide it either.

---

## SEV-2 — S6's "28 oracles" is not a countable pass criterion on this tree

**Anchors:** `0011:S6`, widened to `0011:A4` (Evidence) and
`evidence/grounding/findings.md`.

S6's pass criterion is "all 28 oracles pass, none re-homed under `--all`." A4's
Evidence claims "all 28 `flow next` oracles across
`internal/cli/flow_next_0005_test.go`, `flow_mvv_0005_test.go`,
`flow_adversarial_0005_test.go`, `flow_input_0005_test.go`,
`flow_encoder_0005_test.go`, `flow_resolve_0005_test.go`,
`flow_harness_0005_test.go`, and `reserved_key_0008_test.go` were enumerated and
classified" — but the enumeration itself is not in the record and not in
`evidence/`. I searched: the only files mentioning the count are
`evidence/grounding/findings.md` (which quotes A4 and S6 rather than listing) and
another persona's file. No artifact carries the 28 names.

Counting the tree directly does not recover it. Those eight files hold
`12 + 6 + 5 + ... ` top-level `func Test` declarations across a much larger
population, and only some are `flow next` oracles; `flow_next_0005_test.go`
alone has 12 top-level tests plus 3 `t.Run` subtests. Without the list, "28"
cannot be checked, and neither can "none re-homed" — the meaningful half of the
claim.

This one is partly self-mitigating: S6 itself says "a green 0005 suite is NOT
evidence C1 shipped — S2/S3 are," and the `oracle` mini-check marks S6 an
**anti-oracle**. So S6 is not load-bearing for C1. It is still load-bearing for
the negative control (did the default flip break something unexpected), and that
is what the missing list costs.

**Blocks:** the S6 count assertion and the "MOVES class empty" assertion. The
weaker "the package compiles and `go test ./internal/cli` is green modulo C3's
named re-homings" remains writable, so this is SEV-2 rather than SEV-1.

---

## SEV-3 — `unknown` element ordering is unspecified across a two-source merge

**Anchors:** `0011:C1` (`unknown` clause), `0011:D-identity`, `0011:S6`
("determinism" among the oracles that must keep passing).

Post-change, `unknown` merges facts from two producers: the kernel (which sorts
its undecided atoms via `compareUndecidedAtoms`,
`internal/resolve/guard.go`) and the CLI's own walks for owned keys and gate ids
(insertion order over `row.RequiresOwned` and `row.Gate`,
`internal/cli/flow_next.go:193-207`). No contract states the merged order, and no
scenario asserts it.

`D-identity` requires that "the same request in the same mode over the same
model revision and artifact contents reports the same candidate set" — set, not
sequence — so strictly this is satisfiable. The shipped determinism oracle
(REQ-111, `internal/cli/flow_next_0005_test.go:414-445`) compares two
invocations in one process, which any deterministic merge passes.

So this is a specification gap rather than an unwritable test: a consumer
diffing payloads across builds has no stated guarantee, and I have no criterion
by which to fail a build that reorders `unknown`.

**Blocks:** nothing outright. It prevents a golden-payload/byte-stability oracle
over `unknown`, which the record does not call for. Noted because `§cross-cutting-concerns`
is an unfilled template stub and its "canonical-form / determinism" prompt is
exactly the one that would have caught this.

---

## Not findings — checked and clear

- **S1** is the strongest scenario in the record: `models/rdr.toml` at
  `stage=resolved`, 21→3, named rules, live baseline captured in
  `evidence/spikes/a6-rdr-toml-narrowing.md` with the artifact JSON and the
  observed `outcomes: ["advance","revise","abandon"]`. Writable as-is, and the
  21-vs-3 count is a genuine discriminator against a stripped-match build.
- **S4** is writable: `evidence/spikes/a5-decision-table-no-tag.md` fixture B
  demonstrates `row-recognized-only` with `"unresolved": []` on the live build,
  and `internal/table/normalize.go` lifting `recognized` into `Row.Outcome` is
  confirmed by A2. Empty-match-pattern ⇒ candidate follows from
  `matches(nil) == true` (A3).
- **S2**'s "`--all` rejected by `resolve`/`read-state`/`set-state`" is writable:
  `internal/cli/root.go:44,75` documents the cobra unknown-flag error converter,
  so the exit code and error class are already conventional.
- **C3**'s "4 reads in `flow_next_0005_test.go`" is accurate as a textual count
  — I confirmed occurrences at lines 103, 163, 166, 404. Its `flow_adversarial`
  (446) and `flow_mvv` (66) counts also check out. The count is fine; only the
  gate-id reason (SEV-2 above) blocks the re-homing.
- **C2**'s payload-shape clause (`unknown` as `[]`, not omitted, in both modes)
  gives MVV step 7 a clean pass/fail. Writable.
