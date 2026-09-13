# Deviations — RDR 0021 Lint's normalized-graph export

Phase 2 implementation deviations. Each entry carries exactly one Type and
a line-leading `Status:`.

---

## D1 — REQ-18 has no observable of its own

Type: IMPL-DECISION

Status: mechanical translation

REQ-18 reads: "Exact field spellings are normative AS SPELLED HERE — the
Illustrative Code is an exhibit that disclaims literal assertion and shows
only some members, so it binds nothing."

The clause has no observable of its own, so it carries an EMPTY test column
in `coverage.md` and stays uncovered after Phase 2.

Its two halves are each discharged elsewhere:

- The POSITIVE half — that the field spellings are normative as spelled —
  is REQ-19…REQ-25, which name every member individually and are asserted
  by `graph_document_0021_test.go` against the shipped document.
- The NEGATIVE half — that the Illustrative Code exhibit binds nothing —
  forbids asserting the exhibit. No test in this record asserts it, which
  is what discharges the clause. A test written to cover REQ-18 directly
  would assert the ABSENCE of an assertion, which is vacuous.

Recorded as an implementation decision rather than an author decision: the
record is correct as written and needs no amendment. The demote is a
statement about testability, not about the clause's content.

---

## D2 — REQ-109's test carries two mutually unsatisfiable loops

Type: TEST-FIXTURE

Status: needs author decision

`TestReq109_AngleBracketsAndAmpersandsReachTheWireUnescaped`
(`internal/cli/graph_document_0021_test.go:1040`) cannot pass against any
document, correct or incorrect. Its two loops contradict each other:

```go
for _, escaped := range []string{`<`, `>`, `&`} {
    if strings.Contains(body, escaped) {   // fails when the body HAS `<`
        t.Errorf("the document carries the HTML escape %s; …")
for _, raw := range []string{"a<b", "x&y"} {
    if !strings.Contains(body, raw) {      // fails when the body LACKS `a<b`
```

The fixture authors the tag domain `["a<b", "x&y"]`, so any document
carrying the authored values contains `<` and `&` as literal bytes. Loop 2
requires exactly that; loop 1 refuses exactly that.

The REQ is satisfied by the shipped code. REQ-109 and the XC clause require
the shared non-HTML-escaping encoder so `<`, `>`, `&` "serialize as
THEMSELVES" rather than as `<`/`>`/`&` — which is what loop
1's own failure message says. The export routes through
`clierr.WriteJSONLine` (`marshalGraphJSON`, `internal/cli/graph.go`), and a
live run over the fixture emits `"domain":["a<b","x&y"]` with ZERO `<`
escapes. The intended oracle is the escape SEQUENCES; loop 1 tests the raw
characters instead.

No production change can satisfy both loops. HTML-escaping the document to
clear loop 1 would break loop 2 and violate REQ-43/REQ-109 outright, so the
implementation is deliberately left as the record requires and this entry
carries the defect. Phase 1's tests are not this leg's to amend.

Remedy for the author: loop 1's needles become `<`, `>`, `&`.

---

## D3 — REQ-67's self-scan matches its own source

Type: TEST-FIXTURE

Status: needs author decision

`TestReq67_NoLoadAfterExportInverseIsClaimed`
(`internal/cli/graph_stability_0021_test.go:154`) scans every
`graph_*_0021_test.go` file for a `load ∘ export` inverse, keyed on the
literals `table.Load([]byte(body`, `table.Load([]byte(stdout`, and
`table.Load([]byte(doc`.

Those literals appear in the SCANNING TEST'S OWN SOURCE, at lines 178–180,
as the needles themselves. `graph_stability_0021_test.go` is inside the
scanned set, so the file matches itself and the test fails by construction.

The clause it guards holds: no test in this record feeds an exported
document back to the loader, and the first half of the same test — that the
document carries no authored-model key at top level — passes against the
shipped exporter.

No production change can clear it; the scan's subject is test source.

Remedy for the author: exclude the scanning file from its own sweep, or
split the needles so they do not appear as contiguous literals.

---

## D4 — REQ-85's function splitter misattributes a body

Type: TEST-FIXTURE

Status: needs author decision

`TestReq85_NoTestInThisRecordAssertsAnExitCodeAlone`
(`internal/cli/graph_stability_0021_test.go:195`) requires every test
function that asserts an exit code to carry a byte- or set-equality oracle
beside it, and reports
`graph_surface_0021_test.go::TestReq14_TheVerbUsesOnlyTheExistingZeroTwoExitMapping`.

REQ-14's body genuinely carries `clierr.ExitCodeFor` and none of the
allowed content-oracle needles (`assertRefusalCode`,
`assertNoDocumentOnStdout`, `assertSameSet`, `assertGolden`,
`decodeDocument`, `decodeGeneric`, `ce.Code`, `ce.Param`) — confirmed by
scanning lines 416–442. The test is an exit-mapping table over four arms,
which is what REQ-14 itself specifies, so the constraint and the test it
flags were authored against each other.

No production change can clear it: the subject is this record's own test
source, which Phase 2 does not amend.

Remedy for the author: give REQ-14's arms a content oracle beside the exit
assertion, or exempt a test whose REQ is the exit mapping itself.

---

## D5 — the exported edge relation admitted escape rows the traversal skips

Type: IMPL-GAP

Status: mechanical translation

Fixed in `fc1cbbb`. `internal/graphlint/export.go::edgesFrom` admitted each
satisfiable ESCAPE row as a self-loop, while `successorsOf`
(`reach.go:207-209`) skips every escape row before it asks satisfiability
and `Reach`'s contract defines an edge as a normalized NON-ESCAPE row. The
exported relation was therefore strictly larger than the one lint
certifies, contradicting A2's "EXACTLY the edges the traversal itself took
— equality, not plausibility, is the bar" and `0021:C4`'s structural
no-drift rule. With a catch-all escape row the export minted one phantom
self-loop per reachable node, so the blast radius grew with the model.

The fix moves the escape skip ahead of the satisfiability test, mirroring
`successorsOf` exactly. `reach.go` is untouched (REQ-53/REQ-56 stay
additive) and `ReachGraph` is the only caller. Recorded as IMPL-GAP, not an
author decision: the record was right and the code diverged from it.

---

## D6 — `domain`'s member source is undecided for the finite non-`enum` kinds

Type: SPEC-UNDER

Status: needs author decision

REQ-20 binds `domain` to be "present exactly when `guard.AssignmentCount`
reports it finite", and `internal/cli/graph_document.go:185` gates instead
on `finite && len(decl.Domain) > 0`, contradicting both the REQ and the
function's own doc comment (`:51-54`). The gate half is a genuine
divergence (Phase 3a's FAIL-20).

It cannot be closed mechanically, because clearing the gate requires
deciding what the array CARRIES for the finite kinds that populate no
`decl.Domain`, and no authority spells one:

- `decl.Domain` is non-empty exactly when the kind is `enum` — the loader
  admits `domain` on `enum` alone (`internal/table/load.go:875`) and
  `guard.agrees` requires `!hasDomain` for `bool`, `int`, `set`, `scalar`
  (`internal/guard/declaration.go:206-231`). On the `enum` arm the extra
  conjunct is redundant, so the gate change alone moves no byte.
- `bool` is finite by construction and `int` declares a `{min..max}` bound
  (`0003:C3`, `0003:C6`); a `set` declares an element universe whose
  assignment space is the powerset, which `0003:C24` forbids reading as
  `|universe|`. Those clauses fix per-kind COUNTS, never a wire member
  list, and `guard` exports no member enumerator.
- `0021:C2` states only the presence rule; the Illustrative Code's
  `["draft","final"]` is an `enum` and binds nothing (REQ-18).

Emitting `["false","true"]`, `Elements`, a powerset, or an enumerated bound
would each add a normative wire vocabulary with no REQ-N — ADDITIVE IS NOT
EXEMPT — so this leg records rather than invents it. The committed golden
`internal/cli/testdata/graph_0021_state_machine.json` pins the finite
`bool` tag `flag` with `domain` absent, and `assertGolden` fixes the remedy
for a pin break as "a deliberate schema decision, not a regenerated
golden".

Remedy for the author: decide what `domain` carries for `bool`, `int` and
`set`, or narrow REQ-20's presence rule to the `enum` arm whose members
exist; the gate and the golden then move together.

---
