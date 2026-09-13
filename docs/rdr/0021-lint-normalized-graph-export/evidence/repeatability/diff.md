Model: claude-opus-5[1m]

# Repeatability-lite diff — RDR 0021, run-1 (claude-sonnet-5) vs the record

Neutral diffing context; authored neither the RDR nor run-1. Compared run-1
against `kind=="C"` (C1-C5), `MVV`, `S1-S9`, and the `D-*`/`RT*` elements
first.

## Health signal

**Healthy.** Run-1 carries 12 explicit GUESS markers, and they do not spread
evenly over the record — they cluster on exactly two seams: (a) Go identifiers
for symbols the RDR mandates as capabilities but never spells, and (b) the
`rows[]` member spellings C2 declares normative by pointing at an exhibit that
does not contain them. Everything C2/C3/C5 actually fixes — the `values`
tag-keyed-object shape, the `dot` single-member wrap, the `[]`-not-`null` and
absent-not-null rules, the sort orders, the `<opaque>` sentinel, the
`abstraction` token — run-1 reproduced without a marker and without error.
That is the signature of a record whose determinacy gaps are localized, not
one that is being read confidently-wrong. No rerun on another model needed on
this evidence.

Of the 12 GUESS markers, 8 are discarded below the admissibility bar (naming
and helper-decomposition taste), 1 is a run-1 reading error over a contract the
RDR does fix, and the rest point at the real silences admitted below, which the
GUESS set alone did not fully cover — two admitted findings (F1, F3) are places
run-1 diverged *silently*, with no marker at all, which is the more informative
half of this pass.

## Rejected — GUESS over something the record fixes

- **Edge-recovery mechanism** (run-1 §2 helper 1, "GUESS on internal
  decomposition — whether it re-walks `successorsOf` post-fixpoint"). Not a
  silence. `0021:C4` states the neutrality oracle is MECHANISM-INDEPENDENT and
  "binds equally if Resolve picks A2's in-traversal edge observer (premortem
  P-14)"; `0021:A2` records post-hoc recovery as Verified with the observer as
  the named fallback, and S3 asserts the oracle "holds whichever edge mechanism
  A2 settles on". The openness is a *stated, bounded* choice with an equality
  bar, not an under-determined contract. Rejected.
- **`rows[]` snake_case spellings** (run-1 §3, "GUESS on exact JSON key
  spellings for `source`, `next`, `requires_owned`, `gate`, `escape`, `emit`").
  The *values* are fixed: C2 defers to "RDR 0002's dump field list", and that
  list is a closed, verbatim-pinned vocabulary in code
  (`internal/table/dump.go::dumpColumns` — `identity, source, kind, outcome,
  atoms, next, writes, requires_owned, gate, escape, emit`, fixed "VERBATIM as
  the normative fixtures author it"). Run-1 guessed right by symmetry. The
  reading gap that forced the guess is nonetheless real and is admitted as F2
  below in its accurate form — the defect is C2's self-reference, not the
  spellings.
- **`--emit` CLIError `Param` spelling; `NodeCeiling()` body; every new Go
  identifier** (`newGraphCmd`, `runGraph`, `ReachWithEdges`, `Edge`,
  `BuildDocument`, `RenderDOT`, `export.go`/`dot.go` file names; `table.Load`
  call shape). Naming and helper-decomposition taste. Discarded per the bar.

## Admitted findings

### F1 — `--emit` validation has no position in the arm order

- **Lands on**: `0021:C1` (also `0021:S8`, `§pre-lock-mini-checks` `disposition`)
- **How run-1 rendered it**: run-1's pseudo-code (§4) places the `--emit`
  validity check *after* `os.ReadFile` and after the loader call — a caller
  passing `--emit=xml` against an unreadable path gets `model-unreadable`, and
  against a malformed model gets `model-invalid`. Run-1 marked no GUESS here;
  it presented the placement as fixed, writing "the control flow below is NOT a
  guess — C1 fixes it verbatim against lint.go::runLint's literal arm order".
- **The silence**: C1 does fix an order for the arms it *inherits* — "Model
  selection mirrors `lint`'s arm set ... and codes verbatim
  (`internal/cli/lint.go::runLint`)", and `runLint`'s literal sequence is
  mutual-exclusion → `--flow`-alone → `flag-required` → read → load. But
  `--emit` is new to this verb; `lint` has no such flag, so mirroring supplies
  no position for it. C1 introduces it in a separate, later sentence
  ("`--emit <format>` selects the document ... any other value →
  `flag-invalid-value` naming `emit`") with no ordering clause. S8 lists the
  refusing arms as an unordered set and asserts only "lint's code spellings
  verbatim plus `flag-invalid-value` naming `emit`"; the `disposition` table
  likewise pairs each input class with a code one row at a time, never across
  rows. Nothing in the record says which code wins when two arms are live at
  once. Two implementers get different observable codes for the same argv, and
  S8 passes for both — a flag-validation-first reading is equally conformant
  and arguably more natural, since argument validity normally precedes I/O.
- **Determinacy class**: unstated step order / precedence between refusal arms.

### F2 — C2 declares field spellings normative by citing a partial exhibit

- **Lands on**: `0021:C2` (with `§illustrative-code`)
- **How run-1 rendered it**: run-1 reproduced the `rows[]` keys correctly but
  flagged them GUESS, reasoning that "C2 names these as English labels ... but
  only gives the full worked JSON shape for `identity/kind/outcome/atoms/writes`
  in §illustrative-code's partial example", and fell back to inferring
  snake_case by symmetry with `single_valued`/`requires_owned`.
- **The silence**: C2's closing clause is "Exact field spellings are normative
  as the Illustrative Code spells them", then enumerates
  `rows[…0002's field list…]` — a placeholder rather than a spelling.
  §illustrative-code opens by declaring itself non-binding and incomplete:
  "Illustrative only — tests must not assert these literally. The JSON below is
  a PARTIAL document: the row and tag objects show a few of C2's members, not
  all." So the one clause that makes spellings normative delegates to an exhibit
  that (a) disclaims literal assertion and (b) omits six of the eleven row
  members. The loop closes only by leaving the record entirely, into
  `internal/table/dump.go`. Note this is a *reference* defect, not a value
  defect — the spellings are in fact determinate in code; what is
  under-determined is which artifact an implementer is bound by, and the
  golden-fixture tripwire S2 is built to pin whichever shipped first.
- **Widened past**: C2 → `§illustrative-code` → out of the record into
  `internal/table/dump.go` (`dumpColumns`), since no RDR element carries the
  list.
- **Determinacy class**: unnamed field owner / normative reference that does not
  resolve inside the record.

### F3 — DOT terminal-node marking has no satisfaction predicate

- **Lands on**: `0021:C2` (DOT paragraph; also `0021:A6`, `0021:S6`)
- **How run-1 rendered it**: run-1 carried the requirement through verbatim as
  prose — "initial/terminal nodes marked" (§3) — and its `RenderDOT(doc)
  (string, error)` signature takes only the document. It never states how a
  node is *determined* terminal, and produced no mechanism in pseudo-code. No
  GUESS marker: the gap was reproduced as an unelaborated restatement, which is
  how a silence looks when a reconstructor notices nothing to reconstruct.
- **The silence**: C2 requires "the initial node and terminal-satisfying nodes
  marked", and makes the marking normative ("its node/edge SET and the marker
  are normative, its styling/attributes are not"). But `terminal` on the wire is
  a set of *predicate sets* (`[[{key, operator, literal[], block}]]`), while
  reach nodes are merged, over-approximate value sets whose tags may carry the
  `<opaque>` sentinel. Deciding "node N satisfies terminal predicate set P"
  over a merged node is a real evaluation with at least two live readings — does
  a node whose `values[stage]` is `["draft","final"]` satisfy `stage eq final`
  (some-member) or not (all-members)? — and a third question on top: what an
  `<opaque>` value does to satisfaction, given `0021:A9` fixes only that the
  sentinel reaches the wire, not how it evaluates. A6 asserts derivability
  ("nodes, edges, initial, and terminal-satisfaction are enough") but treats
  terminal-satisfaction as an *input* it is given, never defining it; S6
  deliberately scopes its oracle to the node/edge *identifier* set and the
  marker's presence in the header, so an implementation that marks the wrong
  node set — or no terminal node at all on a merged fixture — passes S6, S3 and
  the MVV unchanged. Whether the marking is even mechanically checked is
  therefore unpinned along with the predicate.
- **Widened past**: C2 and S6 into `0021:A6` and `0021:A9`; no element defines
  the predicate.
- **Determinacy class**: unstated transform / unspecified evaluation semantics
  over merged and abstracted values, with no oracle covering it.

### F4 — the `model` member's value is unfixed

- **Lands on**: `0021:D-identity` (the identity element that would own it; also
  `0021:C2`)
- **How run-1 rendered it**: run-1's `Document` struct carries
  `Model string \`json:"model"\`` with the value shown only as the placeholder
  `"<model identity>"` (§3), unresolved against any source field. No GUESS
  marker — run-1 copied the record's own vagueness forward.
- **The silence**: C2 requires the document carry "model identity and class"
  and fixes the *key* as `model`. `D-identity` fixes identity for the three
  other entities precisely and by source symbol — "a node is its canonical node
  key (`reach.go::(Node).key`) ... a row is RDR 0002's row identity; an edge is
  the `(from, to, rule)` triple" — and stops, never saying what a *model's*
  identity is. `A3` cites `internal/table/model.go` field-walk including
  `Model.ID`, and §illustrative-code shows `"model":"flow"`, both consistent
  with `Model.ID`; but neither is normative (A3 is an assumption about coverage,
  and the exhibit disclaims literal assertion per F2). The competing reading is
  live and consequential: the `--model <path>` argument is the other obvious
  "model identity", and a path-valued member makes the document
  invocation-dependent — the same model exported from two checkouts diffs,
  which is precisely the CI-diffability use the Problem Statement is built on.
  C3's determinism narrowing is to "(model, build)" and would not catch it,
  since a single fixture path is stable within one test run.
- **Widened past**: C2 and D-identity into `§critical-assumptions` (A3) and
  `§illustrative-code`; no normative element fixes the value.
- **Determinacy class**: unnamed field owner / unfixed identity source for a
  required wire member.

## Spans read

`0021:C1`-`C5`, `MVV`, `S1`-`S9`, `D-identity`, `D-wire-byte-format`,
`D-naming`, `D-selection-predicate`, `RT1`-`RT3`, `A1`-`A9`, `F1`-`F5`.
Widened past the contract/scenario spans into `§illustrative-code`,
`§pre-lock-mini-checks` (`authority`/`oracle`/`fidelity`/`disposition`/`trace`),
`§implementation-plan` Phases 1-4, `§capability-dependencies`,
`§existing-infrastructure-audit`, `§approach`, `§problem-statement`,
`§key-discoveries`, `§consequences`, `§risks-and-mitigations`, `§references`.
The Finalization Gate sub-sections (`G-*`) are unfilled templates at this stage
and resolved nothing. Left the record for `internal/table/dump.go` (F2) and
`internal/cli/lint.go::runLint` (F1) to confirm two silences were silences
rather than citations I had not followed.
