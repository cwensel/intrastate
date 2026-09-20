Model: claude-fable-5-1

## Delta scope

Second cove pass, delta-scoped to the elements rewritten while resolving
iter-1's 12 findings: A1, A3, A4, C3, C4, C5, F1, S4, S7, §Approach,
§Research Findings, plus the five new `#### Pre-lock mini-checks` tables
(`authority`, `disposition`, `fidelity`, `oracle`, `trace`). The 12
iter-1 findings are not re-litigated; only new false claims, new
contradictions, and new silences at those anchors are reported. Source
read from `main` (090c539).

## Verification results (source claims added by the rewrite)

All checked against `main`; every claim below held unless marked.

- `internal/cli/flow_next.go::probeRow` (l.469) takes `row table.Row,
  view, owned, observed, all` — no model; its non-test caller (l.259)
  is inside the handler that holds `req.model`. HOLDS (the phrase "takes
  only a `table.Row`" is loose — it takes five args — but the load-
  bearing point, no model at the site, is true).
- `guardSeam()` (flow_resolve.go l.561) is zero-arg; non-test call sites
  are exactly three: flow_resolve.go l.258 (resolve path, `req.model`
  in scope), l.445 inside `escapeClassOf` (probe built from
  `req.model.KernelTable()` at l.430), flow_next.go l.492 inside
  `probeRow`. HOLDS for non-test code; four more call sites exist in
  `escape_shape*_0009_test.go`, which A3's "rg over non-test sources"
  framing excludes. Not filed.
- `Evaluator{}` non-test construction sites: product.go l.542,
  reach.go l.507, flow_resolve.go l.561 — exactly three. HOLDS.
- `valueSatisfies` is called from `Denotation(m *table.Model, …)` one
  frame up (l.437→l.490); `atomAdmitsValue` is two frames below
  `matchSatisfiable` (l.449 → `ownedAtomSatisfiable` l.466 →
  `atomAdmitsValue` l.484/502). HOLDS.
- `matchSatisfiable` keeps `a.Block == table.BlockMatch` (and owned
  provenance) and consults no guard atom; `atomAdmitsValue` constructs
  `guard.Evaluator{}`. HOLDS.
- `guard_mvv_test.go::conformingContractSeam` is a `struct{}` (l.489)
  with a raw-string `eq` arm, driven directly by
  `TestGuardEvaluatorContract` (l.462) and wrapped in `recordingSeam`
  (l.256). HOLDS.
- `guardcontract.go`: one atom construction in the loop with
  `Key: "subject"` (l.51) — every case shares the key. HOLDS.
- `renderWrites` (normalize.go l.580) calls `conform(decl, "eq",
  members)` (l.624) and files kind/domain refusals under
  `CatMalformedTagDeclaration` with the D3 comment. HOLDS.
- `guard.Conforms(m *table.Model, v View)` (declaration.go l.280)
  checks required-key presence and single-valuedness only, never kind;
  no non-test caller. HOLDS.
- `assignment.go::valueAssignments` int arm renders `strconv.Itoa(n)`
  (l.305); called from product.go with `decl := m.Tags[key]`. HOLDS.
- `load.go::valueMembers` (l.1748) renders float64 via
  `FormatFloat(t,'g',-1,64)` (l.1757) and is called before `conform`
  in `atom` (normalize.go l.61→172), `loadInitial` (load.go
  l.1688→1701) and `renderWrites` (l.613→624). Spike run in-tree
  (pelletier/go-toml v2): `eq = -0.0` decodes to float64 `-0`, renders
  `"-0"`, `Atoi("-0")` = 0 nil, `Itoa(0)` = `"0"`; `Atoi` also admits
  `"00"`, `"01"`, `"+1"`. The `fidelity` table's bare-float row HOLDS.
- Refusal categories named in C5 all exist in `internal/table/
  category.go`: `malformed_predicate_atom`, `malformed_initial_
  declaration`, `malformed_tag_declaration`, `emit_value_out_of_domain`;
  CLI codes `flow-tag-invalid` / `flow-write-invalid` in
  `flow_input.go` l.36/40 via `codeInvalidFor`. `canonicalValue`
  (l.738) reaches `conformKind` through `table.ConformValue`. HOLDS —
  but see F-2 for what that implies.
- `atomsFromBlock` doc (normalize.go l.17-19) reads "…the tag-key
  declaration rule are enforced identically in all three blocks";
  refusal `CatUnknownTag` "references the undeclared tag" (l.28). A1
  quote HOLDS.
- `OwnedSnapshot` (accessor/model.go l.424) → `runReaders`
  (flow_exec.go l.332) → `resolve.Resolve` → `assemble`; no conform on
  the path. HOLDS.
- `TagSet.matches` byte-compares (`tv.value != w.Value`, resolve.go
  l.178). HOLDS.
- `kernelCode(KindGuardUnevaluable)` = `flow-guard-unevaluable`;
  `flow next` reason vocabulary is `absent | uncomparable |
  not-evaluated` (flow_next.go l.106); `uncomparable` in `flow next`
  is read off `Refusal.Undecided` (l.561). HOLDS.
- Peer quotes: `0007:A17` headline "nor tag declarations at evaluation
  time", Evidence "The declared kind needed to *parse* the value is
  known to the evaluator from the table it was built for, not from the
  kernel", If-wrong "the seam takes whatever RDR 0003 shows the
  evaluator needs (e.g. the declared kind on the atom) — but never the
  view". `0003:C8` and `0003:C9` as quoted. `JDR 0004 §JD-1` / `§JD-3`
  as cited; JD-3 lands in 0030 Phase 2. All HOLD.

Cross-clause pins (brief item 3):
- C4 vs S4 on the zero-value assertion: same predicate, S4 widens the
  pattern to `var`/`new`/embedded — agree.
- C3 fixture `map[string]string` vs C1 `NewEvaluator(kinds
  map[string]string)` — agree.
- C5 category list vs `disposition` row "per-ingress category (C5's
  list)" — agree by reference; venue wording disagrees (F-2).
- Key Discoveries A17 bullet vs §Approach A17 paragraph — same narrowing
  ("no declaration it was not built over"), both cite the Evidence and
  If-wrong — agree.

Joint-satisfiability of the MVV under C1–C5 (`trace` table): satisfiable
as a design — no clause forbids what another requires — but the table's
own walk is not self-consistent (F-1).

## FINDINGS

### F-1 — `trace` walks the MVV against two different guard literals
- anchor: 0012:§pre-lock-mini-checks (`trace` rows 1, 3, 3b); inherited
  from 0012:MVV step 1 vs step 3; S7 and the `oracle` table
- class: b
- evidence: `trace` step 1 authors "rows `iter eq 3` + fallback"; step 3
  witnesses "parsed `7 == 7` → match" for reader value `07`; step 3b
  witnesses "parsed `4 != 3` → GuardFalse". S7's "same `mvv-int` model"
  declares the guarded row `iter eq 7`; `oracle` row 3 reads "`07` vs
  `eq 7`"; MVV step 1 says `iter eq 3`, MVV step 3 says "against the
  guard `iter eq 7`".
- what: under the guard the table itself authors (`eq 3`), a reader
  value `07` parses to 7 ≠ 3 and the seam answers GuardFalse — the row
  prunes and the parsed-comparison witness does NOT fire, so the step-3
  "consistent" verdict is false as tabled; 3b's `4 != 3` then contradicts
  step 3's `7 == 7`. One literal must be pinned. `7` is the one S7's
  shared fixture and the typed-compare spike use, so MVV step 1 and
  `trace` step 1 should author `iter eq 7`, and 3b should read `4 != 7`.

### F-2 — C5 names the CLI ingress but the record frames the break as load-only
- anchor: 0012:C5; 0012:F1; 0012:§pre-lock-mini-checks (`authority` row
  2, `disposition` row "Authored non-canonical int literal")
- class: c
- evidence: C5's first sentence says non-canonical spellings "are
  REFUSED at load" and its code list ends with "`flow-tag-invalid` /
  `flow-write-invalid` at the CLI (`internal/cli/flow_input.go::
  canonicalValue`)". `canonicalValue` calls `table.ConformValue` →
  `conformKind`, so a round-trip placed in `conformKind` reaches `--tag`
  and `--write` automatically. `authority` row 2 says "load is the only
  spelling authority" while listing `flow_input.go::canonicalValue` as a
  call site; `disposition` says "refused at load". F1 enumerates three
  surfaces (owned reader, model load, `flow next`) and none is the CLI.
  The record's own spike `evidence/spikes/a4-typed-compare.md` (l.65-81)
  measures `--tag iter=07` ACCEPTED today and flipping GuardFalse →
  GuardTrue under typed comparison — i.e. the CLI tag path changes
  behavior either way: refused under C5-in-`conformKind`, or verdict-
  flipped under C2 if C5 stays load-only.
- what: one field pinned two ways. Either C5's venue is "every
  `conform`/`ConformValue` ingress — model load AND the CLI `--tag`/
  `--write` path" (then F1 and the `disposition` table need a fourth
  row: `--tag n=07` / `--write n=+1`, accepted yesterday, now
  `flow-tag-invalid`/`flow-write-invalid` with the canonical rewrite
  named, and `authority` row 2's "load is the only spelling authority"
  becomes "`conformKind` is the only spelling authority"), or C5 is
  load-only and the two CLI codes come out of its list — in which case
  F1 must still name the CLI `07` verdict flip the spike measured. The
  MVV is unaffected (its step 3 is the reader path, which C5 does not
  reach, as `trace` step 3 correctly says).

### F-3 — Failure-mode labels no longer match projector ids after the A1 bullet was inserted
- anchor: 0012:§failure-modes; cited from 0012:S4 ("(F5)"), the
  `disposition` table ("LOUD (F5)"), the `oracle` table ("F4's drift
  scenario")
- class: d
- evidence: `recs inspect 0012` lists `0012:F5 … F4 Suite/fixture drift`
  and `0012:F6 … F5 Zero-value evaluator survives migration`; the new
  second bullet "Undeclared-key guard flip (A1's If-wrong)" is 0012:F2
  and shifted the two labelled bullets by one.
- what: every "F4"/"F5" citation added or kept by the rewrite now
  resolves through the projector to the wrong element (F5 → drift, not
  zero-value survival). Relabel the two bullets F5/F6 (or drop the
  inline numbers, since only these two carry them) and update the three
  citations.

### F-4 — C3's blanket re-keying puts `gte` and `contains` cases on matrix-illegal keys
- anchor: 0012:C3 (second paragraph)
- class: d
- evidence: "the existing cases are re-keyed onto an `enum`/`scalar` key
  and the new legs onto `int`, `bool` and `set` keys". The existing
  suite (`guardcontract.go` l.27-45) carries five `gte` cases and four
  `contains` cases. C2 states the matrix "admits only `int`" for
  `lt/lte/gt/gte` and treats `contains` as a set-arm operator; `0003:C8`
  requires each operator to declare the kinds it accepts.
- what: read literally, the fixture would type `gte` over an
  `enum`/`scalar` key and `contains` over an `enum`/`scalar` key — atoms
  load refuses — and the suite's kind-token coverage claim is then
  carried only by the new legs. Not a contradiction with C2's seam
  behavior (both arms are kind-blind by C2's own text), but the sentence
  should route `gte` cases to the `int` key and `contains` cases to the
  `set` key, leaving `eq`/`in` on `enum`/`scalar`.

## Not filed

- A3's "three call sites" for `guardSeam()` is true of non-test code
  (four more in `*_0009_test.go`); the Evidence's "rg over non-test
  sources" framing carries the qualifier.
- `probeRow` "takes only a `table.Row`" — it takes five parameters; the
  load-bearing claim (no model at the site) is exact.
