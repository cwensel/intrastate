Model: claude-opus-5[1m]

# 3amigo iter-2 Persona 3 — QA / Tester (RDR cli/0011, delta-scoped)

Scope: only the six open entries handed to this pass. Read of the record was
through `rdr inspect`; reads of `internal/resolve` and `internal/cli` tests
were direct.

---

## 1. S5's enumeration of changed oracles — **STILL-OPEN** (one oracle missing)

**Partly closed.** `0011:S5` no longer claims the suite "keeps passing". It now
states outright that "The `internal/resolve` suite does NOT pass unchanged, and
'green' is the wrong pass criterion", enumerates four named oracles with a
per-oracle disposition (retired / re-expected / kept-strengthened), and closes
with "`TestAdv0007_3_MatchNarrowingsHold` and every other `internal/resolve`
oracle MUST keep passing unchanged." The specific inversion defect raised last
iteration is fixed: `TestAdv0007_3_EscapeNotNamingTheConflictedKeyStillRescues`
is now correctly listed as **re-expected**, with the right reasoning — its
ordinary row `A` (`match_conflicted_test.go:47-60`, `gateMatchRow`) carries
`Match: [status=Draft, gate=closed]` over a conflicted `gate`, so under `0011:C1`
row `A` refuses `guard_unevaluable` and the clean escape is never reached.

### Verification of the four listed — all four are real, none is spurious

I swept every `Match` block in `internal/resolve/*_test.go`
(`grep -rn "Match:" internal/resolve/*_test.go` plus every `.Match =` / `.Match
= append` mutation) and classified each against `0011:C1`'s three-valued verdict.
The kernel facts the classification rests on: `resolve.go:174-182`
(`TagSet.matches` folds `!ok || tv.conflicted || tv.value != w.Value` into one
`false`), `resolve.go:472` (ordinary candidate filter) and `resolve.go:614`
(`escapeOrRefuse`'s filter).

- `guard_atoms_test.go:930 TestReq78_MatchPatternStillFoldsAbsenceIntoNonMatch` —
  correctly listed. `row.Match = append(row.Match, Tag{Key: absentKey, …})`
  (line 932), `absentKey = "iterations"` (`guard_fixtures_test.go:204`), a key
  `legalInput` never supplies. Asserts the escapable `no_match` rescue. Inverts
  under C1. **Real.**
- `match_conflicted_test.go:…DuplicateKeyMakesRowSelectionAFunctionOfSlicePosition`
  — asserts `got.Refusal.Kind != resolve.KindNoMatch` is an error. Conflicted
  `gate`. Refusal kind moves to `guard_unevaluable`. **Real.**
- `match_conflicted_test.go:…ConflictedKeyIsNotEscapableByNamingIt` — same
  `KindNoMatch` assertion; escape row's own `Match` names the conflicted `gate`,
  so under C1 it is disqualified as a rescuer and the ordinary row's
  `guard_unevaluable` stands. **Real.**
- `match_conflicted_test.go:…EscapeNotNamingTheConflictedKeyStillRescues` —
  **Real**, as reasoned above.

`TestAdv0007_3_MatchNarrowingsHold` is correctly excluded: all three subtests
(`"closed","closed"` same-value repeat; owned-over-observed precedence;
non-duplicate permutation) leave `gate` unconflicted and present-and-equal, so
the verdict stays `match`. Verified against `merge`'s conflict rule
(`resolve.go:235-243`: `next.conflicted = prior.conflicted || prior.value !=
t.Value`, within one provenance only).

### The gap: a fifth oracle changes disposition and is not listed

`internal/resolve/reserved_key_0008_test.go:131
TestReq3_AbsentOutcomeYieldsAViewWithNoRecognizedKey`.

```go
in.Recognized = ""
in.Table.Rows[0].Match = []resolve.Tag{
    {Key: "status", Value: "Draft"},
    {Key: reservedKey, Value: ""},        // reservedKey = "recognized"
}
...
if got.Refusal == nil || got.Refusal.Kind != resolve.KindNoMatch {
    t.Errorf("disposition = %+v; want no_match — the key is absent from the view", got)
}
```

`assemble` binds `recognized` only for a non-empty `Input.Recognized`
(`resolve.go:222-227`: `if in.Recognized != "" { view.tags[recognizedTagKey] = … }`),
which the constant's own doc comment restates (`resolve.go:131-132`, "assemble
injects only for a non-empty Input.Recognized"). So with `Recognized = ""` the
key `recognized` is **ABSENT from the view** — the test's own error string says
so. Under `0011:C1` an absent match key is `indeterminate`, the row is not
selected, and it refuses `guard_unevaluable`, not the `KindNoMatch` this oracle
pins. Currently green (`go test ./internal/resolve/ -run
TestReq3_AbsentOutcomeYieldsAViewWithNoRecognizedKey` → PASS).

This is not a cosmetic omission. It is the second oracle in the package that
pins absence→`no_match`, it lives in a different file from the four listed
(`reserved_key_0008_test.go`, which `0011:A7`'s Method line does not even name
among the files to enumerate), and it belongs to a *different* RDR's suite
(0008's reserved-key contract). It will go red on implementation, and S5's
current text makes "any OTHER oracle going red is a defect" — so a conforming
build hits a stated defect condition on a correct implementation. That is
exactly the discrimination S5 says it provides, pointed the wrong way.

Note also `0011:A7`'s Evidence line lists the files to enumerate as
`match_conflicted_test.go`, `guard_atoms_test.go`, `escape_shape_0009_test.go`,
`escape_shape_mvv_0009_test.go`, `guard_fixtures_test.go` — `reserved_key_0008_test.go`
and `reserved_key_sameview_0008_test.go` are absent from that list, so the
Stage-6 method as written would not have found this one either.

Checked and confirmed NOT changing (so the count is 5, not more):
`TestReq4And5_EmptyOutcomeIsNotClosedByThisRDR`
(`reserved_key_0008_test.go:164`) also sets `Recognized = ""`, but
`!in.Table.models(in.Recognized)` short-circuits to `unmodeled_outcome` at
`resolve.go:458` *before* the candidate loop, so it is unaffected;
`TestAdv0009_AConformingTableKeepsItsDisposition`
(`escape_shape_adv_0009_test.go:110`) uses `stage=elsewhere` — present and
unequal, a true `no-match`; `dormantBreachInput`'s
`Match: status=NeverHeldValue` (`escape_shape_fixtures_0009_test.go:84`) is
present-and-unequal; `TestReq41`'s `Row.Match` mutation
(`reserved_key_0008_test.go:445`) only ever runs through `CheckInput`, never
`Resolve`'s disposition. All other `Match` blocks in the package name
`status` / `reviews` / `recognized` / `stage`, every one present and
unconflicted in its fixture view.

**Blocks:** the S5 pass criterion is unsatisfiable as written — a correct
implementation turns `reserved_key_0008_test.go:131` red, which S5 classifies as
a defect. Either add it as the fifth re-expected oracle (refusal kind
`no_match` → `guard_unevaluable`, its "the key is absent from the view" point
preserved and arguably sharpened), or state why 0008's reserved-key suite is out
of scope. `0011:A7`'s "exactly four … and no others" and its file enumeration
list must move with it, as must MVV step 8's "no OTHER oracle in that package
changes disposition".

---

## 2. `uncomparable` reporting case writable — **CLOSED**

Named surface for the token to arrive on: `0011:C1` states the probe's
`Refusal.Undecided` payload "is the CARRIER: the CLI reads the match facts off
it exactly as it reads guard facts today, and no new field on `Plan` is
required", and `0011:A8` names the two production sites that must widen
(`internal/resolve/guard.go::isGuardBlock`'s filter and the non-`GuardUnevaluable`
payload discard). The `authority` mini-check's census row for "candidate
`unknown` list" names the consuming site concretely: `flow_next.go::summarize`
at `flow_next.go:168`, folding "kernel indeterminate atoms (read off
`Refusal.Undecided`, reasons `absent`/`uncomparable`)". An oracle can be written
against a named field on a named struct read at a named function.

Producer reachability is closed and, better, fenced: `0011:C3` and `0011:S3` both
require the conflicted fixture be built "from two invoked readers writing the
same owned key with differing values", and both name the anti-pattern with its
reason — "a repeated `--tag` cannot reach the kernel
(`flow_input.go::parseTags` refuses it as `flow-tag-duplicate`), so a fixture
written that way tests the CLI's duplicate refusal instead of the match verdict
and MUST NOT be accepted as this fixture." Verified reachable in source:
`internal/cli/flow_exec.go:223` does `owned = append(owned, result.OwnedSnapshot()...)`
per reader with no dedup, and the kernel marks a within-provenance differing
repeat conflicted at `resolve.go:239`. A two-reader fixture therefore lands a
conflicted owned key in the view. Writable.

---

## 3. Un-run gate id reason token — **CLOSED**

`0011:C1` mints it explicitly and bounds it: a gate id un-run because
`--evaluate-gates` was not passed "is neither `absent` … nor `uncomparable` …
so it carries the reason `not-evaluated`, which this RDR mints for the gate
class ONLY and which never appears on a match or owned-key entry."
`0011:D-undecided-vocabulary` carries the same split with its rationale
(`0007:C8` types a refusal payload; `unknown` is a CLI reporting surface with a
third source), and the `disposition` mini-check table has the row
`gate id, --evaluate-gates not passed | yes | {gate-id, not-evaluated} | no | 0`.

This does make `0011:C3`'s re-homing mechanical. C3 now says outright "The
gate-id reads re-home to the `not-evaluated` reason C1 names; a re-homing that
has to invent a reason token is not mechanical and is a defect" — the exact gap
from iteration 1, named and closed with a defined token. The affected read is
`internal/cli/flow_next_0005_test.go:163-171`
(`TestReq41And46And57_WithoutTheFlagGateIdsAreUnresolvedFactsNotResults`,
`slices.Contains(unresolved, "approval")`), which re-homes to
`{approval, not-evaluated}` with no invention required.

C3's count of six also verifies against source. Distinct test reads of the
`unresolved` key: `flow_next_0005_test.go:103, 163, 166, 404` (4),
`flow_adversarial_0005_test.go:446` (1), `flow_mvv_0005_test.go:66` (1) = 6,
matching C3's "(4, at the presence check and the gate-id assertions), (1), (1)".
Minor and non-blocking: of the four in `flow_next_0005_test.go`, line 404 is the
absent-`flag` assertion rather than a gate-id one, so C3's parenthetical
undercounts the classes present; the count itself is right and every site
re-homes to a defined token.

---

## 4. The uncheckable "28 oracles" criterion in S6 — **CLOSED**

`0011:S6` now demotes the number explicitly: "The '28 oracles' figure A4 reports
is a classification count, not a checkable pass criterion: the enumeration
behind it was never written down, so the build cannot verify '28' or 'none
re-homed' against a list." It then offers a two-armed remedy — "(a) emit the
enumeration as an artifact when it re-runs the classification, or (b) assert the
weaker checkable form — `go test ./internal/cli` green except C3's six named
re-homings, with no test re-homed under `--all`" — and closes with "Treat the
number as provenance for A4's verdict, not as an oracle." Arm (b) is executable
as written against a named package and a named six-item list. The anti-oracle
("A green 0005 suite is NOT evidence C1 shipped — S2/S3 are") is retained and is
also carried in the `oracle` mini-check's S6 row.

---

## 5. `unknown` element ordering — **CLOSED**

`0011:D-undecided-reporting-shape` specifies it: "Ordering: entries are sorted
by `(key, reason)`, so the merge of kernel-supplied and CLI-computed facts is
stable across builds — `D-identity` requires only set equality, but an unstable
order would churn any golden payload a consumer diffs." A total order over a
two-field tuple whose second component is drawn from a closed three-member
vocabulary (`absent`, `uncomparable`, `not-evaluated`) is fully determined, so a
stable-output oracle is writable: assert the emitted list equals the
`(key, reason)`-sorted expected list, byte for byte. The dedup rule that governs
what is in the list before sorting is likewise pinned in `0011:C1` ("Entries are
deduplicated by `{key, reason}` pair, not by key, so a key that is both an
undecided match atom and an unestablished owned key reports once per distinct
reason") and restated in the `authority` mini-check census row.

---

## 6. New untestable claims introduced by the rewrite — **CLOSED**

Checked the three named areas plus MVV step 6b.

**C1's site-class split.** The split is stated as a rule with a named
disposition per class, not as a principle: reporting site (`flow next`) MUST NOT
exclude; selection site (`flow resolve`) MUST NOT select, and instead "refuses
`guard_unevaluable` … carrying its undecided atoms in `Refusal.Undecided`". The
selection-site paragraph goes further and pins the two sub-cases separately — an
ordinary row must not be selected *and* "MUST NOT count toward the
`len(selected)` that yields `ambiguous_match` against a row that genuinely
matched", and an escape edge whose own match is indeterminate "is likewise NOT a
viable rescuer". Both are directly assertable, and both correspond to real call
sites I confirmed: `resolve.go:472` (ordinary) and `resolve.go:614`
(`escapeOrRefuse`). The `len(selected)` clause is the one an implementer could
otherwise get wrong by filtering after the count; it is now written.
`0011:D-verdict-use-by-site-class` carries the same content with both rejected
alternatives named. Nothing here is a claim without an oracle.

**A8.** Testable as stated. It names both production sites by file and
construct, states the exact invariant to check — "leaving every existing guard
payload byte-identical" — and names the mechanism that could break it
(`compareUndecidedAtoms` orders by block, so "a new block value participates in
ordering"). Byte-identity of existing guard payloads is a directly executable
oracle: run the shipped guard-path payload assertions before and after the
widening. Its `If wrong` names a concrete design consequence (a separate
`Refusal` field, with `0011:C1`'s carrier clause and
`D-undecided-reporting-shape` moving with it). Status `Pending` with Method
"Source Search + MVV Test" is the honest state; the assumption is verifiable,
not unfalsifiable.

**MVV step 6b.** "Assert the selection-site half in `internal/resolve`: the same
undecidable match atom refuses `guard_unevaluable` (not `no_match`, not a plan)
and names the atom in `Refusal.Undecided`; an escape edge whose own match is
indeterminate does not rescue." Three assertions, each with a named expected
value and two named excluded values. Package-local, and the fixture shape it
needs already exists in-package (`match_conflicted_test.go::gateInput` builds a
conflicted `gate` through `Observed`). Writable as-is.

No new untestable claim found. The one new defect this pass raises (entry 1) is
an incompleteness in an enumeration, not an unfalsifiable claim.

---

## Summary

| # | Entry | Status |
| --- | --- | --- |
| 1 | S5's four-oracle enumeration | **STILL-OPEN** — a fifth changed oracle, `reserved_key_0008_test.go:131`, is missing; A7 and MVV 8 inherit it |
| 2 | `uncomparable` reporting writable | CLOSED |
| 3 | `not-evaluated` gate-id token / C3 mechanical | CLOSED |
| 4 | S6 "28 oracles" | CLOSED |
| 5 | `unknown` ordering | CLOSED |
| 6 | New untestable claims (C1 split, A8, MVV 6b) | CLOSED |

Five of six fixes hold. One blocking gap remains, and it is a completeness
defect in exactly the list the last pass asked to be made complete.
