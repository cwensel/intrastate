Model: claude-opus-5[1m]

# Repeatability DIFF — RDR 0030, full variant (profile: foundational)

Runs compared: run-1 `claude-opus-5[1m]`, run-2 `claude-fable-5-1`,
run-3 `claude-sonnet-5`. Three distinct base models, so a split along the
model boundary is a genuine signal rather than sampling noise.

Findings are ordered 3-way → 2-way → GUESS cluster. Each names the RDR
element the rewrite lands on.

---

## 1. Disagreements

### D1 (3-way) — Where the admits filter, the zero-cell refusal and the bound check RUN
**Element: `0030:C1` (expansion paragraph) — residually `§phase-2-expand-into-cells`.**

| Run | Rendering |
| --- | --- |
| run-1 (opus) | Admits filter, zero-cell refusal, checked add and the C2 domain bound ALL inside `renderWrites`; `expand` receives an already-resolved `(Cells, Literals)` pair and only mints rows. |
| run-2 (fable) | Same placement as run-1 — but the run reaches it by citing `§phase-2`, which it read explicitly. |
| run-3 (sonnet) | `admittedCells(...)`, the zero-cell refusal, `checkedAdd` and the domain bound are all in `expand`; `renderWrites` only records a `StepSpec{N}`. |

C1 is the divergence source. Its expansion paragraph says only "the loader
expands a rule carrying a step write into literal rows", then names
`renderWrites` solely as the site whose SIGNATURE widens — never as the
site that OWNS the cell walk. The split is fixed, but only in
`§phase-2-expand-into-cells` ("That work splits across two existing sites,
because `expand` cannot do it alone … The declaration-dependent half …
stays in `renderWrites`"), which is an Implementation Plan section, not a
contract. A reader held to the Normative Contracts cannot recover it: C1
even places the `int` width check, the `#`/duplicate guards and the
step/clear collision in prose that reads as loader-generic. run-3's
placement is a conformant reading of C1 alone.

Consequence of run-3's reading, and why this is algorithmic rather than
decompositional: `expand` has no `*loader` receiver and therefore no
`TagDecl`, so run-3 had to invent a `TagDecl` reach that does not exist,
AND had to widen `expand`'s return to carry an error (see D2). The
under-determination is the OWNER of a step, not the name of a helper.

Rewrite: hoist §phase-2's two-site split into C1 — name `renderWrites` as
the owner of admitted-cell evaluation, the width check, the per-cell
literal, the `conform` bound and the member guards, and `expand` as the
row minter that receives the resolved per-key cell/literal pair.

### D2 (3-way) — `expand`'s return signature: does it refuse?
**Element: `0030:C1` (the two-signature-widenings paragraph).**

| Run | Rendering |
| --- | --- |
| run-1 | `expand(base Row, predicates []Atom, outcome Atom, writes []writeSpec) []Row` — no error; every refusal is `renderWrites`'. |
| run-2 | `expand(base Row, predicates []Atom, outcome Atom, writes []WriteSpec) []Row` — no error. |
| run-3 | `func expand(writes []TagValue, predicates []Atom) (rows []Row, err error)` — a NEW error return, and `expand` refuses `malformed_tag_declaration` at three points. |

C1's own words are "widening that pair is one of the TWO signature changes
this clause implies … Both widenings are local to `normalize.go` and
neither adds a site." C1 enumerates exactly two widenings — the `writes`
parameter type and `renderWrites`' new `predicates` parameter — and never
says `expand`'s RETURN is untouched. `§phase-2` quotes the current
signature `expand(base Row, predicates []Atom, outcome Atom, writes
[]TagValue) []Row`, but only inside an argument about the `*loader`
receiver, and C1 is silent on whether a third widening (an error return)
is admitted. run-3 minted an unstated error mode on a function C1 treats
as total. This is a signature the RDR leaves open.

Rewrite: C1 should state that `expand` remains total (`[]Row`, no error)
and that every refusal this record mints fires before `expand` is called.

### D3 (3-way) — The carrier the step spec rides in
**Element: `0030:C1` ("the spec itself rides alongside as one more per-key record").**

| Run | Rendering |
| --- | --- |
| run-1 | A NEW type `writeSpec{Key, Literal, Step *stepSpec}` REPLACES `[]TagValue` in both signatures; `stepSpec{Cells []string, Literals []string}` carries two parallel slices. |
| run-2 | A NEW type `WriteSpec{Key, Literal, Cells []cellWrite}` replaces `[]TagValue`; `cellWrite{Cell, Literal}` pairs each cell with its literal — a slice of pairs, not parallel slices. |
| run-3 | `TagValue` itself is WIDENED IN PLACE with an optional `Step *StepSpec` field; the signatures keep `[]TagValue`. `StepSpec` carries only `{N int}` — NOT the resolved cells, because run-3 resolves cells later, in `expand`. |

Three structurally distinct data models for the same crossing. C1 says the
pair "cannot carry a non-literal today — widening that pair is one of the
TWO signature changes", which admits both "widen `TagValue`" and "replace
`TagValue`" readings, and never states WHAT the record carries: the
unresolved step magnitude (run-3) or the resolved per-cell literals
(run-1, run-2). That second question is not naming — it determines whether
the cell walk happens before or after the record is built, so it is the
same under-determination as D1 seen from the data side, and it is why
run-3's `StepSpec{N}` is internally consistent with its D1 placement.

Rewrite: C1 should state that the record crossing `renderWrites` → `expand`
carries the RESOLVED admitted cells and the rendered literal per cell
(one entry per cell), so `expand` needs neither the `TagDecl` nor the step
magnitude.

### D4 (2-way, splits opus+fable vs sonnet) — Where the step/clear collision is checked
**Element: `0030:C1` (clear-collision paragraph) / `0030:A14`.**

- run-1: the collision check runs AFTER the per-cell walk, and run-1 flags
  its own placement as a GUESS, noting the record fixes only that the spec
  be DERIVED before the clear pass.
- run-2: the check runs BEFORE the member enumeration (`if key in
  src.clear: refuse` immediately after the kind gate), and separately GUESSes
  whether `renderWrites` inspects `src.clear` directly or receives the list.
- run-3: the check runs after the `#`/duplicate guards but before any cell
  walk — same ordinal slot as run-2.

C1 fixes ONE ordering ("The step spec is therefore derived BEFORE the
clear-list pass, so the pair is detectable") and A14 adds a different
hint ("the step spec is available at the clear loop, so the pair can be
refused THERE rather than needing a second pass"). These do not agree on
which loop emits the refusal: C1 implies the write loop holds the spec and
the clear pass detects, A14 implies the clear loop refuses. That matters
for the OBSERVABLE first-failure ordering when a rule both steps a key
into an out-of-domain cell AND clears it — run-1 reports the C2 bound,
run-2/run-3 report the clear collision. The RDR fixes a first-failure
order for cells (C2: "the FIRST in the tag's own domain order") but fixes
no order BETWEEN the collision refusal and the bound refusal.

Rewrite: C1 should state which refusal wins when a key is both stepped
past its domain and named in `clear`.

### D5 (2-way, splits fable vs opus+sonnet) — Ownership of the `int` width check vs the enum member guards relative to the cell walk
**Element: `0030:C1` (grammar/kind arm).**

- run-1 and run-3: `resolve.IntWidth` / `intWidth` is consulted on the `int`
  arm ONLY, and the `#`/duplicate guards on the `enum` arm only, both
  before cells are enumerated — the kind switch owns both.
- run-2: same placement, but run-2 folds the width call INTO
  `domainMembers(decl)` (`int: width,ok := resolve.IntWidth(min,max); !ok →
  refuse`), making the width refusal a consequence of enumeration rather
  than a separate gate.

Marginal, and admitted only because C1 reasons about the width check as
the thing that decides whether the loader "can count the cell set, let
alone walk it" — a precedence claim — while never stating whether the
width refusal precedes or is produced by the enumeration. Under run-2's
folding, a declaration that is both unrepresentable AND carries a `#`
member reports the width failure; under run-1/run-3's the same. No divergent
output was produced by the three runs, so this is the weakest admitted
finding; recorded because the precedence is genuinely unstated, not
because the runs disagreed on an observable.

### D6 (2-way, splits opus vs fable+sonnet) — The lint read-back join key
**Element: `0030:C3` (the JOIN KEY paragraph).**

- run-1: commits to `Span` as THE join key, and flags the choice as a GUESS.
- run-2: "joins on `Span` or the recovered bare id" — leaves both live and
  posits a recovery helper `authoredID(identity string) string`.
- run-3: "stays the bare `Span` (or recovered authored id)" — also leaves
  both live.

C3 writes "the in-package read-back joins on a key that does not carry it
— the row's `Span` … **or** the authored id recovered from the suffixed
form", and the `authority` mini-check table repeats the disjunction
("the in-package join key is `Span` (or the recovered authored id)"). The
RDR states a disjunction where a single key is required: the two are not
interchangeable for `groupHasOverlap`, which matches against
`ruleIDsOf(g)` — a list of bare `RuleID`s, NOT spans — so joining on
`Span` requires ALSO changing what `ruleIDsOf` returns, while recovering
the authored id does not. The RDR leaves the field owner unnamed and
the two arms have different blast radii.

Rewrite: C3 should pick one and say which, or state that `ruleIDsOf`
changes with the `Span` arm.

---

## 2. GUESS clusters (two or more runs marked GUESS)

### G1 — The widened per-key write record's type and fields (3/3 GUESS)
**Element: `0030:C1`.** run-1 "GUESS (shape)" on `writeSpec`/`stepSpec`;
run-2 "GUESS: exact shape" on `WriteSpec`/`cellWrite`; run-3 "GUESS at
fields (RDR does not give the struct literally)" on the widened `TagValue`
+ `StepSpec`. The strongest cluster in the set, and it is D3's other face:
all three flagged it, and all three produced a different model. The rewrite
is the same as D3's.

### G2 — The extracted admits shim's signature and name (3/3 GUESS)
**Element: `0030:C1` / `§phase-2` / `authority` mini-check row 1.**
run-1 `AtomAdmits(atom GuardAtom, value string) GuardVerdict` — GUESS on
name only, taking the string form from §prerequisites. run-2
`AtomSatisfies(atom GuardAtom, held Value) Verdict` — GUESS on name AND on
whether the verdict type is `GuardTrue/GuardFalse/GuardUndecided` or another
spelling. run-3 does not name it at all, listing only `intWidth` and
`IntDomain` under the shim and leaving the admits shim's Go surface blank.

Beyond naming — which is inadmissible — the admitted part is the HELD-VALUE
TYPE and the VERDICT TYPE. §prerequisites fixes the held value as the cell
STRING ("the shim's admits call passes the cell string, NOT the declared
kind"), but the Normative Contracts never carry it, and run-2 rendered it as
an abstract `Value` while run-3 rendered nothing. The verdict type is fixed
nowhere: C1 says "three-valued" and the `authority` table names the
callers' collapses (`!= GuardFalse`), but no element states the type.

Rewrite: C1 (or a new surface of C1) should carry the shim's signature —
`(resolve.GuardAtom, string) → <the three-valued verdict type>` — since
the held-value type is load-bearing for 0012's later typing work and today
lives only in §prerequisites.

### G3 — `intWidth`'s exported spelling after the move (2/3 GUESS)
**Element: `0030:C1` (the `intWidth` relocation paragraph).**
run-1 "**GUESS:** it is exported as `resolve.IntWidth` since `internal/guard`
and `internal/table` both call it across a package boundary." run-2 "GUESS:
exported spelling; RDR names it ::intWidth and says it relocates as-is."
run-3 writes it as lowercase `intWidth` in `internal/resolve` with no GUESS
marker — i.e. run-3 rendered an UNCALLABLE function, since `table` and
`guard` both must reach it across the package boundary.

Admissible because it is not naming taste: C1 says the core "relocates
as-is", and as-is it is unexported, which contradicts C1's own requirement
that the loader and lint share it. The RDR states a relocation whose
visibility it does not resolve, and one run produced an incoherent result
from the literal reading.

Rewrite: C1 should say the moved core is EXPORTED, since its two callers
are now in different packages.

### G4 — Whether `internal/resolve` pre-exists (1 GUESS, but a 2-way factual split)
**Element: `0030:C1` / `§phase-2`.**
run-3 explicitly GUESSes "whether `internal/resolve` is a brand-new package
or already exists holding `Evaluator` pre-record", concluding it "may be
newly created by this work". run-1 and run-2 both treat it as pre-existing
and merely receiving the move. The record settles it — A4 cites
`internal/resolve/resolve.go::GuardEvaluator` and
`internal/resolve/guardcontract.go` as existing, and §phase-2 says
"`resolve` already declares `GuardEvaluator`" — but C1 calls it "the
`internal/resolve` shim", a phrase that reads as minting. A `§phase-2`
reader recovers it; a contracts-only reader does not.

Low-severity; recorded because run-3 marked it and the phrasing is the cause.

---

## 3. Agreement

All three runs rendered identically, with no GUESS between them: the TOML
authoring surface (`{ step = <n> }`, one key, non-zero integer, admitted on
bounded `int` and non-empty-`domain` `enum` only); the full refusal-category
split (`malformed_tag_declaration` on the write-block path,
`malformed_initial_declaration` under `[initial]`, `malformed_predicate_atom`
in predicate position, none merged); `Failure`'s six fields unchanged with
all four bound-failure parts on `Detail` and the cell named first; the
both-carriers rule (`Row.Writes` AND `Row.NextTags` populated per row, never
aliased); the always-append suffix with its own constant-true `suffixed`
local; subsumption of a same-tag `match in` into the conjunction with a
per-row rewrite to `match eq = <cell>`; key-alone ordering of step points
among all candidates; the exclude-on-undecided collapse at the loader's call
site; first-failure reporting in the tag's own domain order; the checked add
ordered before the render; `IntDomain` staying in `guard`; and C3's
row-names-suffixed / rule-names-bare split with no new vocabulary member.
That is the determinate core of the record, and it is large.
