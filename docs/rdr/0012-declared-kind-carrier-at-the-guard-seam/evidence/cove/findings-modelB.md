Model: claude-sonnet-5

## Step 0 — Grounding

Edge sweep: `/tmp/0012b-edges.json` — 73 edges total (45 source-anchor, 14
jdr, 6 mentions, 4 surface-of, 2 predecessor, 1 joint-decision-home, 1
artifact). **All 45 source-anchor edges resolved=true** — no false/absent
edges. This is NOT an empty starting set, so the projection is trusted and
the sweep below widens past it per instructions (the record's real
`path::Symbol` claims plus its no-symbol claims).

CONFIRMED (compact — read on `main`):
- `internal/guard/grammar.go::Evaluator` — `struct{}`, `eq` arm is
  `value == atom.Literal` (raw string), zero-value form retained today.
- `internal/table/load.go::conformKind` — int arm is bare `strconv.Atoi`,
  admits `"00"`/`"01"`/`"+1"`; `valueMembers`'s float64 arm is
  `strconv.FormatFloat(t, 'g', -1, 64)` before any kind check.
- `internal/cli/flow_input.go::canonicalValue` — calls `table.ConformValue`
  (→ `conformKind`+`conformDomain`); error text matches S7's fixture
  byte-for-byte.
- `internal/accessor/model.go::ReadResult.OwnedSnapshot` — builds
  `resolve.Tag{Key,Value}` from reader bytes, no kind check.
- `internal/resolve/resolve.go::assemble`/`merge` — no kind check on the
  owned-merge path.
- `internal/guard/guard_evaluator_0003_test.go::TestReq34…` — `reflect...
  NumField() != 0` assertion, exactly the "structural proxy" A2 claims.
- `internal/cli/flow_resolve.go::guardSeam`, `internal/guard/product.go::
  valueSatisfies`, `internal/graphlint/reach.go::atomAdmitsValue` — exactly
  three non-test `Evaluator{}` sites on `main` (`rg "Evaluator{}"`
  confirms); each has `*table.Model` one or two frames up exactly as A3
  describes (`guardSeam()` callers hold `req.model`; `Denotation(m
  *table.Model,...)` calls `valueSatisfies`; `matchSatisfiable(m
  *table.Model,...)` → `ownedAtomSatisfiable` → `atomAdmitsValue`).
- `internal/guard/assignment.go::valueAssignments` int arm —
  `strconv.Itoa(n)` renders the domain (A4 HELD-side claim) — see F-1,
  the RDR's own edge/Evidence cites this at `product.go`, wrong file.
- `internal/table/model.go::KernelTable` / `internal/resolve` non-test
  sources — zero `internal/table` imports; A5's one-way direction holds.
- `internal/table/normalize.go::atomsFromBlock` — docstring quote "the
  tag-key declaration rule are enforced identically in all three blocks"
  is verbatim; `CatUnknownTag`/"references the undeclared tag" confirmed.
- `internal/guard/declaration.go::DeclarationOf` — sibling kind-lookup
  reused by design; confirmed as the only kind-lookup path (see inverse
  check below).
- `internal/resolve/resolve.go::TagSet.matches` — plain `tv.value !=
  w.Value`; MATCH atoms stay byte-compared, confirming C1's resolution-path
  scope and the "lint/runtime share the seam for GUARD atoms only" claim.
- `internal/resolve/resolve.go::KindGuardUnevaluable` /
  `internal/resolve/guard.go::ReasonUncomparable` — both constants exist
  verbatim, matching C1/C2/F1's refusal-code claims.
- `072c7a0` — commit exists, message matches the Background's inline-fix
  claim; `parseHeldSet`'s nil-check is present at `grammar.go:179-196`.
- `docs/jdr/0004-declaration-carrying-package-and-undecidable-value-venue.md`
  §JD-1/§JD-2/§JD-3 — all three entries exist and bind 0012, matching the
  RDR's citations (JD-3 in particular: "loader-side admitted-cell shim
  constructs through `NewEvaluator`").
- `internal/table/normalize.go::renderWrites` at line 580 — EXISTS on
  `main` today (contra a naive reading of "no fourth site"); it does not
  construct any `Evaluator`, so C4's claim ("constructs no evaluator on
  main") holds, but see F-2 for a wording risk this creates.
- Five-kind vocabulary `{enum,bool,int,set,scalar}` — verbatim in
  `grammar.go:54`; `eq`/`in` matrix rows exclude `set` exactly as C2
  states.

REFUTED / NOT-FOUND:
- **NOT-FOUND (path attribution)**: A4's Evidence cites
  `internal/guard/product.go::valueAssignments` for the HELD-side
  `strconv.Itoa` render. `valueAssignments` is actually DEFINED in
  `internal/guard/assignment.go:273`; `product.go` only CALLS it (4 call
  sites: lines 194, 336, 450, 477). The substantive claim (Itoa rendering)
  is correct; the `path::Symbol` is wrong. See F-1.

Inverse-sibling check (new discriminator: declared-kind lookup at the
value seam): `rg '\.Kind\b'` across the repo shows `decl.Kind`/`.Kind`
consulted in ~20 files, all downstream of `internal/guard/declaration.go::
DeclarationOf` or `table.TagDecl.Kind` directly (load-time conformance,
lint rendering). No second, independent kind-inference mechanism exists
(e.g., inferring kind from literal shape) — this matches the RDR's own
"Briefly Rejected: Kind inference from literal bytes" and its "searched,
no second kind-lookup path exists" claim. CONFIRMED as stated.

Spans widened past: the anchor set covers A1-A5's Evidence fields and a
scatter of inline path citations in Problem Statement / Normative
Contracts / Decision Rationale / Research Findings / Testing Strategy.
I widened past those into: (1) the Finalization Gate section (still raw
template placeholders — expected at this stage, since COVE runs pre-lock
and gate.md is written at Finalize; not a finding), (2) the full
Alternatives/Briefly-Rejected prose (no anchors, pure argument — checked
for internal consistency against the Decision Rationale matrix), (3) the
full Testing Strategy / MVV (no anchors on most scenario bullets — cross-
checked scenario 6/7/8's literal fixture text against actual error
strings), (4) `internal/table/normalize.go` around `renderWrites` and the
three write/initial/match ingress sites, none of which the anchor set
named directly but which C5's "every authoring site" claim requires.

## Step 1-2 — Twelve questions, independent answers

**Q1 (codebase).** Does `internal/guard/grammar.go::Evaluator.Evaluate`'s
`eq` arm today compare raw strings, or does it already parse by kind?
**A:** Raw strings — `case "eq": return boolResult(value == atom.Literal)`,
confirmed by direct read of `grammar.go:106-108`. No kind parsing exists
today. Matches the RDR's premise exactly.

**Q2 (codebase).** Is the "exactly three non-test `Evaluator{}` sites"
claim (A3/C4) still true, or has a fourth site landed since propose?
**A:** True — `rg -n "Evaluator\{\}" --type go | grep -v _test.go` returns
exactly three: `flow_resolve.go:561`, `product.go:542`,
`reach.go:507`. No fourth site on `main` as of this read.

**Q3 (codebase).** Does `internal/table/normalize.go::renderWrites`
(JDR 0004 §JD-3's named future construction site) already construct a
`guard.Evaluator` on `main` today, making C4's "constructs no evaluator on
main" claim stale?
**A:** No — `renderWrites` (line 580) calls `conform(decl,"eq",members)`
for kind/domain conformance but never touches `guard.Evaluator` or
`NewEvaluator`. The claim holds on today's `main`. (Note: the function
itself already EXISTS, unlike what a literal reading of "That site
constructs no evaluator on main" might suggest to a reader who assumes the
whole shim is future work — see F-2.)

**Q4 (codebase).** Does `internal/table/load.go::conformKind`'s int arm
actually accept `"00"`, `"01"`, `"+1"`, `"-0"` as A4/C5 claim?
**A:** Yes — the int arm is exactly `if _, err := strconv.Atoi(member);
err != nil { ... }`; Go's `strconv.Atoi` accepts all four spellings
(verified against Go's documented `Atoi`/`ParseInt` grammar, which permits
an optional sign and leading zeros). Confirmed.

**Q5 (codebase).** Is the switch in the illustrative `eq` sketch
exhaustive over the five-kind vocabulary, or does it silently drop a
kind?
**A:** The sketch (`§Illustrative Code`) has cases for `int`, `bool`,
`enum`/`scalar` (combined), and `default` (undeclared). It never lists
`set` explicitly — `set` falls into `default` and returns
`GuardUnevaluable`, which is exactly what C2's normative text separately
states ("kind `set`: the operator/kind matrix does not admit eq/in over
set ... answers GuardUnevaluable"). So the sketch is exhaustive by
construction (5 kinds → 3 explicit arms + 1 default arm covering both
`set` and undeclared), consistent with the normative prose. No gap.

**Q6 (codebase).** Does a sibling path already make the "declared-kind
carrier at the value seam" decision this RDR proposes — i.e., is there
prior art in this exact repo for typing a comparison by declared kind?
**A:** Yes, partially: `internal/guard/declaration.go::DeclarationOf` and
lint's `product.go`/`assignment.go` already type-render the HELD side
(the domain) by declared kind, and `internal/table/load.go::conformKind`
already type-checks the LITERAL side at load. What does NOT exist today
is a carrier that types the COMPARISON at the value seam itself — the gap
this RDR closes. So the RDR's own "Existing Infrastructure Audit" row 1
(reuse `DeclarationOf`) is accurate, and no independent competing
implementation of the same decision exists to cite as a rejected
duplicate.

**Q7 (codebase).** Does `internal/resolve/resolve.go::TagSet.matches`
already use typed comparison for MATCH atoms, making C1's "resolution
path continues to read no declarations" claim about GUARD atoms
under-scoped (i.e., does it apply to MATCH too, by accident)?
**A:** No — `matches` uses `tv.value != w.Value`, a plain string
inequality; MATCH atoms are untouched by this RDR and remain raw-string
compared post-implementation, exactly as C1 states ("the rule governs
that path" = resolution path generally, and this RDR's typed-comparison
change is scoped to GUARD atoms via the `Evaluator` seam only, never
`matches`). Confirmed no accidental scope creep.

**Q8 (RDR-internal).** Does the Decision Rationale's scored matrix
internally contradict the Load-Bearing Decisions' "C5 is not Alternative
3" note — i.e., does C5 (load-time int canonicalization) actually
reintroduce the "boundary validation" pattern the matrix scores O4 as
losing on?
**A:** RDR is explicit on this and not silent: "C5 constrains authored
model text at load, a venue that already conforms kinds (conformKind) and
whose population is disjoint from the runtime's (JDR 0004 §JD-2). It
tightens an existing load check rather than adding a runtime validation
layer, so the seam still carries its own honesty." This directly
addresses the apparent tension by distinguishing load-time (existing
venue, disjoint population) from runtime caller validation (O4, rejected).
No internal contradiction — the record pre-empts the question in its own
text.

**Q9 (RDR-internal).** Is the "If wrong" field non-empty and substantive
for every Critical Assumption (a Finalization Gate requirement), and does
any assumption's "If wrong" text implicitly promise a behavior this
record's Failure Modes section doesn't also name?
**A:** All five (A1-A5) carry non-empty "If wrong" text. A1's If-wrong
("a guard that decided before this change now refuses
flow-guard-unevaluable — loud, but a behavior change A1 did not predict")
is not separately named in §Failure Modes (F1-F5); F1/F2/F3/F5 cover the
owned-ingress misroute and construction-site misses, but not an
undeclared-key-guard-atom flip specifically. RDR is silent on this
specific sub-case in Failure Modes. Low-stakes (A1 is Verified so the
scenario is "if wrong," a defensive-arm-only path) but it is a gap — see
F-3.

**Q10 (RDR-internal).** Does the Joint-check (JC1, fired → RDR 0030)
correctly identify why 0012 and 0030 collide on the same kind vocabulary,
and is that collision actually resolved (not just noted) by lock time?
**A:** §Decision Rationale states "Its re-fire opened JDR 0004 (settled
2026-09-20, binds 0012 and 0030)" and JD-1/JD-2/JD-3 are all `decided` in
the JDR file (confirmed above). So the joint-check is not merely noted —
it is resolved via JDR 0004, which is dated the same day as the git log's
most recent commits. Confirmed resolved, not dangling.

**Q11 (codebase).** Does the S6 fixture's claimed refusal message ("00"
is not the canonical spelling of int tag n; write 0") match any existing
error-string convention in the codebase, or is it inventing new wording
this RDR hasn't grounded against `conformKind`'s actual error path?
**A:** `conformKind` TODAY returns `"%s is not an int"` (quoted). C5's
proposed diagnostic ("is not the canonical spelling of int tag n; write
0") is NEW text this RDR must add — `conformKind` does not emit it today.
This is expected (C5 is an unimplemented normative contract, not a
description of current behavior) and the RDR is explicit that this is
prospective ("Zero committed models are affected... so the break is
prospective"). Not a defect — the RDR correctly frames C5 as new,
unshipped diagnostic text, distinct from A1-A5 which describe existing
behavior.

**Q12 (codebase).** Is `resolve.GuardEvaluator` really a single-method
interface (as C1 states, "a single-method INTERFACE... not a func type"),
confirming `NewEvaluator`'s constructor-based carrier is even structurally
possible without changing the interface?
**A:** RDR 0007:C1 (peer record, read above) states verbatim: "a
single-method INTERFACE named `GuardEvaluator`, not a func type." I did
not re-grep `resolve.GuardEvaluator`'s Go declaration directly in this
pass (peer-RDR citation accepted per Peer RDR method), but the resolve
package source at `internal/resolve/resolve.go` and `guard.go` is
consistent with this (no func-type declarations were seen while reading
`assemble`/`matches`/refusal-kind constants in that package). RDR is not
silent; cites 0007:C1 directly. Treated as CONFIRMED via peer-RDR
citation, consistent with everything else read in `internal/resolve`.

## Step 3 — Findings

### F-1 — A4's HELD-side evidence cites the wrong file for `valueAssignments`
- anchor: 0012:A4
- class: a
- evidence: `internal/guard/assignment.go:273` (`func valueAssignments(d
  table.TagDecl) ([]string, bool)`) vs. the RDR's cited
  `internal/guard/product.go::valueAssignments`
- what: A4's Evidence field cites `internal/guard/product.go::
  valueAssignments` as the HELD-side `strconv.Itoa` renderer. The function
  is actually DEFINED in `internal/guard/assignment.go`; `product.go` only
  calls it (4 call sites confirmed: lines 194, 336, 450, 477). The
  substantive behavioral claim (int domain rendered via `strconv.Itoa`,
  canonical spelling only) is correct and independently verified at
  `assignment.go:308`. This is a `path::Symbol` misattribution in a
  Verified, load-bearing assumption's Evidence field — exactly what the
  Finalization Gate's Assumption Verification item asks the gate-writer to
  confirm ("each cited `path::Symbol` resolves on `main`"). It does not
  resolve to the named file.

### F-2 — C4's "constructs no evaluator on main" phrasing could mislead a future reader about `renderWrites`'s current existence
- anchor: 0012:C4
- class: c
- evidence: `internal/table/normalize.go:580` (`func (l *loader)
  renderWrites(rule *sourceRule, id string, isEscape bool) (...)`)
- what: C4's dependency clause says "the loader-side admitted-cell shim
  `internal/table/normalize.go::renderWrites`... That site constructs no
  evaluator on `main` — it becomes a construction site only when RDR
  0030's Phase 2 shim extraction lands." This is TRUE as written (no
  `Evaluator` construction happens there today) but the phrasing "becomes
  a construction site... when [0030] lands" could be misread as "the
  function itself doesn't exist yet." `renderWrites` already exists today
  and already performs write-block kind/domain conformance via
  `conform(decl,"eq",members)` (confirmed at line 624). This is not a
  factual error — the claim as literally stated is correct — but it is an
  under-specified claim that a reader verifying "no fourth site" against
  `rg "Evaluator{}"` could momentarily conflate with "no such function."
  Recommend the gate-writer or a future refine pass make explicit that
  `renderWrites` exists today and only its evaluator-construction
  obligation is deferred.

### F-3 — Failure Modes is silent on the A1 If-wrong scenario (undeclared-key guard flip)
- anchor: 0012:F1
- class: d
- evidence: 0012:A1 "If wrong" — "a guard that decided before this change
  now refuses flow-guard-unevaluable — loud, but a behavior change A1 did
  not predict"
- what: A1's If-wrong names a specific behavior change (a previously-
  deciding guard over an undeclared key now refuses loudly) that is
  distinct from F1 (owned-ingress misroute → refusal, the RDR's primary
  visible break) and F5 (missed construction site → refusal). None of
  F1-F5 names the undeclared-key-guard scenario specifically, even though
  A1 flags it as a possible unpredicted behavior change if A1 turns out
  wrong. Since A1 is Verified (not Pending), this is low-probability, but
  the Finalization Gate's own Assumption Verification instructions ask
  whether "If wrong" text is honored elsewhere; this is a minor
  completeness gap in Failure Modes rather than a defect in A1 itself.

### F-4 — S7's fixture error string was independently verified byte-for-byte
- anchor: 0012:S7
- class: (informational — not a defect; recorded to close the loop, no
  action needed)
- evidence: `internal/cli/flow_input.go:751-753`
  (`table.ConformValue`/`userErr` message construction)
- what: Verified the exact wire string S7 pins
  (`"the value for \`iter\` does not conform to its declaration: \"many\"
  is not an int"`) matches the current code's message construction
  exactly, including `strconv.Quote` behavior on the literal. Not filed as
  a defect finding — included because it was the highest-risk "exact
  string" claim in the record and is fully CONFIRMED, closing out a
  candidate false-positive.
