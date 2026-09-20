Model: claude-opus-5

# CoVe pre-lock lens — cli/0030 (findings-b)

Scope note. `§decision-rationale` carries `Ground-sweep: clean (33 anchors)`,
so Step 0 is scoped to claims added or edited AFTER propose (`c00f01e`):
the refine commits `6f42adb`, `586ecba`, the resolve commit `61f4cea` and the
rulings absorb `408722a`. The post-propose additions are substantial — the
whole `internal/resolve` shim-extraction decision, the A4 third-arm fence, the
A8 refutation, the A7 spike, and Phase 2's repoint — so the sweep below covers
every symbol those commits introduced plus the inverse (sibling-already-decides)
check the lens requires.

Edge worklist from `recs inspect --json --filter edges,elements 0030`: 94 edges,
93 `resolved=true`, 1 `ABSENT` (the `artifact` edge to `{ARTIFACT_DIR}/gate.md`,
an unexpanded template token, not a source anchor). The empty NOT-FOUND set is
NOT a clean record — the projector resolves a `path::Symbol` edge by path, not by
symbol, which is exactly how F1 below slipped through. I read the record whole.

Spans widened past: (1) the edge worklist, because it was empty of `false`/absent
entries — I re-verified every post-propose `path::Symbol` by grep regardless;
(2) `Row.Writes` (the only write carrier C1/C3 name) out to `Row.NextTags`,
because `internal/table/normalize.go:834-837` populates both from the same
`writes` slice and six shipped consumers read `NextTags`; (3) `TagDecl.Domain`'s
`#` ban out to its duplicate/emptiness checks, because C1's enum step is
position-indexed over that slice.

## Step 0 — grounding ledger

| Claim (RDR) | Verdict | Cite |
| --- | --- | --- |
| `internal/table/normalize.go::normalizeAtom`'s `case "in":` arm is the `#` ban site (A3 Evidence) | **NOT-FOUND** | no `normalizeAtom` anywhere in the repo; the arm lives in `func (l *loader) atom(...)` at `internal/table/normalize.go:43`, ban at `:120-128` |
| The `#` ban has exactly three sites; `tagDecl` performs none on `src.Domain` (A3) | CONFIRMED | `suffixSep` guards at `normalize.go:124` (`in` members), `normalize.go:297` (rule ids), `load.go:163` (outcome alphabet); `load.go:848-905` `tagDecl` has no `#` check |
| `tagDecl` carries `Domain: src.Domain` unsorted (A2) | CONFIRMED | `internal/table/load.go:899` `Domain: src.Domain` |
| `renderWrites` runs BEFORE `expand`; expand's tail clones one write set onto every row (A1) | CONFIRMED | `normalize.go:444` then `:471 return expand(base, predicates, outcomeAtom, writes)`; tail `:836-837` |
| The choice-point discriminant is `c.atom.Operator == "in" && c.atom.Block == BlockMatch` (A1) | CONFIRMED | `normalize.go:123` |
| `compareAtoms` is the choice-point sort (A1) | CONFIRMED | `normalize.go:266-277`, used at `:746-748` |
| `Row.KernelRow` drops `Suffix` (A8) | CONFIRMED | `internal/table/model.go:362` builds `resolve.Row{RuleID: r.RuleID, …}`; no `Suffix` field |
| `flow_resolve.go::rowByID` first-matches on bare `RuleID` (A8, audit) | CONFIRMED | `internal/cli/flow_resolve.go:408-415` |
| `flow_next.go::summarize` publishes the bare rule id (A8) | CONFIRMED | `internal/cli/flow_next.go:337` `Rule: row.RuleID` |
| `graph_document.go` publishes `identity` per row (A8, C3) | CONFIRMED | `internal/cli/graph_document.go:87` `Identity string \`json:"identity"\``, set at `:214` `r.Identity()` |
| `dump.go` publishes an `identity` column (A8, C3) | CONFIRMED | `internal/table/dump.go:18,114` |
| `internal/guard/grammar.go::Evaluator.Evaluate` is three-valued; `valueSatisfies` returns `resolve.GuardResult` (A4) | CONFIRMED | `internal/guard/product.go:541-553`, doc comment "it returns the seam's THREE-VALUED verdict" |
| The operator/kind matrix confines `lt/lte/gt/gte` to `int` (A4, C1) | CONFIRMED | `internal/table/model.go:110-119` `operatorKinds` |
| `conform` kind-checks an ordered atom's bound (A4) | CONFIRMED | `internal/table/load.go:1834-1838` — `conformKind` only, NOT `conformDomain` |
| `conform`/`conformDomain` is the literal-write conformance check C2 reuses | CONFIRMED | `internal/table/normalize.go:624` `conform(decl, "eq", members)` inside `renderWrites`; `conformDomain` at `internal/table/load.go:1869-1892` |
| `valueMembers` admits string/bool/int64/float64/[]any and refuses the rest as "is not a tag value" | CONFIRMED | `internal/table/load.go:1748-1783` |
| `go list -deps ./internal/table` names neither `guard` nor `graphlint` (A7) | CONFIRMED | output is exactly `internal/resolve`, `internal/table` |
| `internal/guard/declaration.go::AssignmentCount`, `::IntDomain`, `::intWidth` exist | CONFIRMED | `declaration.go:80`, `:241`, `:190` |
| `intWidth` reports a full-width span as unusable (F1) | CONFIRMED | `declaration.go:190-200`, `span == math.MaxInt` → `false` |
| `checkIdempotentWrites` is match-block `eq` only, so a step row never trips it (audit) | CONFIRMED | `internal/graphlint/groups.go:352-356` `if atom.Block != table.BlockMatch \|\| atom.Operator != "eq" { continue }` |
| `internal/graphlint/taxonomy.go::CodeOverlap` / `CodeIdempotentWrite` exist (A5) | CONFIRMED | `taxonomy.go:49-55,112` |
| `TestReq92_…` and `TestReq34_…` exist in `internal/cli/decision_table_0010_test.go` | CONFIRMED | `:859`, `:888` |
| `valueSatisfies`, `atomAdmitsValue`, `renderSet`, `renderSetLiteral` exist where cited | CONFIRMED | `guard/product.go:541`, `graphlint/reach.go:502`, `guard/assignment.go:353`, `graphlint/reach.go:519` |
| "the set renderer twice" — `renderSet` and `renderSetLiteral` are two copies of one function (audit) | **REFUTED** | `renderSet` marshals `members` as authored; `renderSetLiteral` first calls `canonicalValues` (`reach.go:343-345`, `slices.Compact(slices.Sorted(...))`). Different functions |
| "it is already duplicated three ways … C1's 'as the runtime evaluator would' then holds structurally" (audit / Phase 2) | **REFUTED in part** | `graphlint/reach.go:502-514` collapses the third arm the OPPOSITE way from C1: `return verdict != resolve.GuardFalse` (admit on UNEVALUABLE, documented "the false-green direction"). It is not a copy of `valueSatisfies`; it is `valueSatisfies` plus an admit-on-unknown collapse |
| `Row.Writes` and `Row.Atoms` alone are what post-normalization surfaces read (C3) | **REFUTED** | `Row.NextTags` is a separate field populated independently (`model.go:293-295`, `0002:C15`) and read by `graphlint/reach.go:327` (`writesOf`), `graphlint/engine.go:147` (`Fingerprint`), `cli/graph_document.go:219`, `cli/flow_next.go:397`, `guard/lint.go:183`, `table/dump.go:125` |
| "no existing X bans `#` in a domain member" (A3, inverse check) | CONFIRMED as a genuine gap | see the three-site row above |
| Inverse: does a sibling already decide which domain members a rule's atoms admit? | **sibling exists** | `internal/guard/product.go::Denotation` (`:437`) + `::acceptedIn` (`:591`) and `internal/graphlint/reach.go::atomAdmitsValue` (`:502`) each already make this decision; the RDR cites both and declines reuse, which is disclosed — but the third-arm mismatch (above) is not |
| `[initial]` refuses an inline table under the same category as a write (C1's interception arm) | **REFUTED (category)** | `internal/table/load.go:1690` files it as `CatMalformedInitialDeclaration`, not `CatMalformedTagDeclaration` |

## Step 1 — questions

Codebase-only (must be answered from source):

1. Does `internal/table/normalize.go` contain a symbol named `normalizeAtom`?
2. Is `Row.NextTags` populated separately from `Row.Writes`, and does any shipped
   consumer read `NextTags` without `Writes`?
3. Is `internal/graphlint/reach.go::atomAdmitsValue` a behaviour-identical copy of
   `internal/guard/product.go::valueSatisfies`, such that repointing it at one
   extracted shim is behaviour-preserving?
4. Are `guard/assignment.go::renderSet` and `graphlint/reach.go::renderSetLiteral`
   the same function?
5. Does `internal/table/load.go::tagDecl` reject a repeated member in an enum
   `domain`, or an empty `domain` on an `enum` kind?
6. Does `internal/table/load.go::conform` apply `conformDomain` to an ordered
   (`lt`/`lte`/`gt`/`gte`) atom's bound?
7. Does `internal/table/normalize.go::expand` emit a suffix element for a
   single-member choice point today?
8. Does `internal/graphlint/engine.go::Fingerprint` read `Row.Writes`,
   `Row.NextTags`, or both?
9. Does `go list -deps ./internal/table` reach `internal/guard` or
   `internal/graphlint`?
10. Does `internal/graphlint/groups.go::checkIdempotentWrites` consult
    `guard.all` atoms?
11. Does `[initial]` refuse a non-tag-value under the same CLIError category the
    write path uses?
12. (RDR-internal) Does C1's grammar arm state which side of the `Writes`/
    `NextTags` pair the per-cell stepped literal lands on?

## Step 2 — answers

Each answered independently against source; no answer leans on another.

**A1 — No.** `grep -rn normalizeAtom internal/` returns nothing. The `#`-ban arm
A3's Evidence attributes to `normalize.go::normalizeAtom` is inside
`func (l *loader) atom(decl TagDecl, key, operator string, raw any, b Block, owner string)`
at `internal/table/normalize.go:43`; the guard itself is at `:120-128`.

**A2 — Yes and yes.** `internal/table/model.go:293-295`: "NextTags and Writes are
populated independently from the rule; one is never an alias of the other."
`internal/table/normalize.go:836-837` sets both from `cloneTagValues(writes)`.
`internal/graphlint/engine.go:147` reads `canonicalTags(row.NextTags)` and never
`row.Writes`. `internal/cli/graph_document.go:219` publishes
`Next: graphTagValues(r.NextTags)`.

**A3 — No.** `internal/graphlint/reach.go:502-514`:
`func atomAdmitsValue(a table.Atom, held string) bool { … return verdict != resolve.GuardFalse }`,
with the doc comment "An UNEVALUABLE verdict is taken as satisfiable: the
relation over-approximates, and pruning an edge lint cannot decide is the
false-green direction." `valueSatisfies` returns the three-valued
`resolve.GuardResult` unchanged. The two differ by a deliberate, documented
third-arm collapse, and that collapse is the OPPOSITE of C1's ("an atom answering
neither true nor false at a member does NOT admit that member").

**A4 — No.** `guard/assignment.go:353-360` `json.Marshal(members)` on the members
as given. `graphlint/reach.go:519-531` first calls `canonicalValues(a.Literal)`
(`reach.go:343-345` = `slices.Compact(slices.Sorted(...))`) and substitutes `[]`
for nil. Sorted-and-compacted vs. as-authored.

**A5 — No to both.** `internal/table/load.go:848-905` `tagDecl` checks
provenance, kind token, kind/field agreement, the single-valued marker and
`min > max`. There is no `firstDuplicate` call on `src.Domain` (the only two
`firstDuplicate` sites are `normalize.go:113` and `:139`, both on atom members),
and no `len(src.Domain) == 0` check for `enum`.

**A6 — No.** `internal/table/load.go:1834-1838`:
`case "lt", "lte", "gt", "gte": … return conformKind(decl, members[0])`.
`conformDomain` is reached only by the fall-through loop at `:1841-1848`.

**A7 — No.** `internal/table/normalize.go:131-132`:
`suffixed := expanding && len(members) > 1` — "A suffix element is emitted only
for an atom with MORE THAN one member, so a single-member `in` mints no second
identity."

**A8 — `NextTags` only.** `internal/graphlint/engine.go:147`
`for _, t := range canonicalTags(row.NextTags)`. (`writesOf` at
`graphlint/reach.go:322-334` is the one place that unions the two.)

**A9 — No.** `go list -deps ./internal/table | grep intrastate` prints exactly
`internal/resolve` and `internal/table`.

**A10 — No.** `internal/graphlint/groups.go:352-354`:
`if atom.Block != table.BlockMatch || atom.Operator != "eq" { continue }`.

**A11 — No.** `internal/table/load.go:1690`:
`return fail(CatMalformedInitialDeclaration, "[initial] "+key+": "+err.Error())`,
against `CatMalformedTagDeclaration` on the write path (`load.go:614,624,631`).

**A12 — RDR is silent.** C1 says "the literal write `cell + n` (`int`) or the
domain member `n` positions from the cell (`enum`)" and closes "`Row.Writes`
holds literals only". C3 says the surfaces "read expanded rows through
`Row.Writes` and `Row.Atoms` alone". Neither `NextTags` nor `0002:C15` appears
anywhere in the Normative Contracts fence; the only two occurrences of `NextTags`
in the record are narrative (`§problem-statement`, `§alternatives-considered`).

## Findings

### F1 — A3's Evidence cites a symbol that does not exist
- anchor: 0030:A3
- class: not-found
- claim: A3 Evidence — "the one that bans a member is `internal/table/normalize.go::normalizeAtom`'s `case \"in\":` arm, gated `if b == BlockMatch`".
- evidence: `grep -rn "normalizeAtom" internal/` returns nothing. The arm is in `func (l *loader) atom(...)`, `internal/table/normalize.go:43`, guard at `:120-128`: `if b == BlockMatch { for _, m := range members { if strings.Contains(m, suffixSep) { … } } }`.
- why it matters: the claim A3 makes is TRUE (verified independently at the three real sites), but the Evidence line names a non-existent symbol, so the Gate's "each cited `path::Symbol` resolves on `main`" check fails, and the projector's `resolved=true` did not catch it — it resolves by path. Re-point A3 at `internal/table/normalize.go::loader.atom`.

### F2 — C3's invariance enumerates `Row.Writes` but not `Row.NextTags`
- anchor: 0030:C3
- class: silence
- claim: C3 — "`flow next`, `flow resolve`, `flow set-state`, `lint`, `dump`, and `graph` read expanded rows through `Row.Writes` and `Row.Atoms` alone". C1 closes with "`Row.Writes` holds literals only".
- evidence: `internal/table/model.go:293-295` — "NextTags and Writes are populated independently from the rule; one is never an alias of the other" (`0002:C15`). `internal/table/normalize.go:836-837` sets `rows[i].NextTags = cloneTagValues(writes)` AND `rows[i].Writes = cloneTagValues(writes)`. `NextTags`-only readers: `internal/graphlint/engine.go:147` (`Fingerprint`), `internal/cli/graph_document.go:219` (`Next`), `internal/cli/flow_next.go:397`, `internal/guard/lint.go:183`, `internal/table/dump.go:125`, `internal/graphlint/reach.go:327`.
- why it matters: this is the load-bearing silence. If Phase 2 writes the per-cell stepped literal into `Writes` only — which is literally what C1 specifies — then every expanded row of one stepped rule keeps the IDENTICAL `NextTags`, so `Fingerprint` (NextTags-only) collides across all N rows, `graph`'s `next` publishes the unstepped value, and `reach.go::successor` traverses to the wrong node. C1/C3 must say the stepped literal replaces the key in BOTH carriers, per row, and must cite `0002:C15` for why saying it once is not enough.

### F3 — `graphlint::atomAdmitsValue` is not a copy of `valueSatisfies`; repointing it is not behaviour-preserving
- anchor: 0030:§existing-infrastructure-audit (the `valueSatisfies` row) / Phase 2
- class: refuted
- claim: audit — "It is already duplicated three ways (`guard/product.go::valueSatisfies`, `graphlint/reach.go::atomAdmitsValue`, a test-local copy in `internal/resolve`) … repoint `guard` and `graphlint`"; Phase 2 — "repoint `guard::valueSatisfies` and `graphlint::atomAdmitsValue` at it."
- evidence: `internal/graphlint/reach.go:502-514` — `func atomAdmitsValue(a table.Atom, held string) bool { … return verdict != resolve.GuardFalse }`, doc comment: "An UNEVALUABLE verdict is taken as satisfiable: the relation over-approximates, and pruning an edge lint cannot decide is the false-green direction." Against C1: "an atom answering neither true nor false at a member does NOT admit that member."
- why it matters: the two live callers collapse the third arm in OPPOSITE directions, and the loader (this record's new caller) takes graphlint's opposite. Phase 2's "repoint" is therefore a three-caller change where each caller must retain its own collapse at the call site, and the extraction is behaviour-preserving only under that condition. The audit's "Reuse (extract)" cell and "C1's 'as the runtime evaluator would' then holds structurally" both overstate: what holds structurally is the two-valued core, never the third arm — which is precisely the arm C1 spends a paragraph on.

### F4 — "the set renderer twice" is two different functions
- anchor: 0030:§existing-infrastructure-audit (the `valueSatisfies` row)
- class: refuted
- claim: audit — "plus the set renderer twice (`guard/assignment.go::renderSet`, `graphlint/reach.go::renderSetLiteral`)".
- evidence: `internal/guard/assignment.go:353` `renderSet` = `json.Marshal(members)` on the sequence as given. `internal/graphlint/reach.go:519` `renderSetLiteral` = `json.Marshal(canonicalValues(members))`, `canonicalValues` at `:343-345` being `slices.Compact(slices.Sorted(slices.Values(value)))`, plus a nil→`[]` substitution.
- why it matters: an extraction that folds the two "duplicates" into one renderer silently changes one of the two `in`-literal comparisons. Either the loader's admits filter gets a member-order-sensitive `in` (if it takes `renderSet`) or graphlint's edge relation loses its canonicalization (if it takes the other). The audit should say "two renderers, not two copies", and Phase 2 should name which one the loader takes.

### F5 — C1's enum step is position-indexed over a domain that admits duplicates and admits emptiness
- anchor: 0030:C1
- class: silence
- claim: C1 — "the domain member `n` positions from the cell (`enum`)"; and "admitted only on a tag whose declared kind is … `enum` with a non-empty `domain`".
- evidence: `internal/table/load.go:848-905` `tagDecl` — the only `Domain` checks are `len(src.Domain) > 0 && src.Kind != "enum"` (`:876`) and the carry `Domain: src.Domain` (`:899`). No duplicate check (the repo's `firstDuplicate`, `normalize.go:200`, is called only at `:113` and `:139` on atom members) and no empty-domain check for `enum`. `conformDomain` (`:1872`) is guarded `if len(decl.Domain) > 0`, so an empty enum domain conforms every member today.
- why it matters: `domain = ["a", "b", "a"]` is authorable on `main`, and "n positions from the cell" has two answers for cell `a`. C1 names a `#`-member refusal it needed A3 to discover; the duplicate-member case is the same class of discovery and is unstated. Either C1 refuses a step over a domain with a repeated member, or D-identity must say which occurrence indexes (and note that the suffix `retry#a` is then ambiguous across two cells — the same collision `#` was banned to prevent).

### F6 — C1's "`[initial]` keeps the literal-only grammar" does not name the refusal category, which differs from the write path's
- anchor: 0030:C1 (and 0030:S4)
- class: refuted
- claim: C1 — "`[initial]` values and predicate literals keep the literal-only grammar; the table shape is admitted on the write-block path alone." S4 — "`[initial] attempt = { step = 1 }` and a predicate literal `{ step = 1 }` both refuse."
- evidence: `internal/table/load.go:1688-1690` — `members, err := valueMembers(l.doc.Initial[key]); if err != nil { return fail(CatMalformedInitialDeclaration, "[initial] "+key+": "+err.Error()) }`. The write path files under `CatMalformedTagDeclaration` (`internal/table/normalize.go:615,625,638`).
- why it matters: C1's whole grammar arm is stated in terms of `malformed_tag_declaration`, and a reader takes "keeps the literal-only grammar" to mean "refuses the same way". It refuses under a different code, and S4 asserts only "refuse" with no code, so the MVV would pass whichever code fires. One clause fixes it: the `[initial]` refusal is `malformed_initial_declaration`, the write-path one `malformed_tag_declaration`.

### F7 — D-identity reverses the suffix rule `expand` implements, and states it as a parenthetical rather than a contract clause
- anchor: 0030:D-identity
- class: contradiction
- claim: D-identity — "a single-cell expansion still appends its element (a stepped row is never the authored row, so it always carries the cell — the opposite of the single-member `in` rule, whose one row IS the authored row)."
- evidence: `internal/table/normalize.go:131-132` — `suffixed := expanding && len(members) > 1`, commented "A suffix element is emitted only for an atom with MORE THAN one member, so a single-member `in` mints no second identity." The step point must therefore set `suffixed` unconditionally, i.e. the flag stops being a function of `len(members)` and becomes a function of the choice-point KIND.
- why it matters: C1's expansion paragraph says only "the cell appended to its expansion suffix" and is silent on the single-cell case; the reversal lives in a Load-Bearing Decision parenthetical. `expand`'s `suffixed` is one boolean shared by every choice point, so this is a real code fork, and the identity-tuple totality argument `0002:C13` rests on (`normalize.go:104-113`) is the thing being varied. It belongs in C1's fence, where a peer can cite it.

### F8 — A4's Evidence proves the third arm unreachable using a check that does not run on ordered atoms
- anchor: 0030:A4
- class: refuted
- claim: A4 Evidence — "`internal/table/load.go::conform` kind-checks their bound, so both `Atoi` calls succeed over an `IntDomain()` member".
- evidence: `internal/table/load.go:1834-1838` — the ordered-operator arm returns `conformKind(decl, members[0])` and never reaches `conformDomain`. `conformKind` (`:1853-1865`) is `strconv.Atoi` for `int`, nothing more.
- why it matters: the claim as written is CORRECT and the conclusion holds — `Atoi` is exactly what `conformKind` does, so both parses succeed. But the reader who checks it will find `conform` has two arms and that the ordered one skips the domain check entirely, which is a live consequence this record does not state: `lt = 500` on `min=0,max=9` loads today, so a step rule may author an ordered atom whose bound is outside the domain. That admits every cell (harmless here — the conjunction is evaluated per member) but it means "the rule's own positive atoms" can be vacuous in a way `conform` never reports. Worth one clause in C1 or a line in A4 narrowing the Evidence to `conformKind`.

### F9 — the `unless`-over-admits mitigation names a refusal that fires only when the bound is crossed
- anchor: 0030:C1 (the `unless` paragraph) / 0030:§risks-and-mitigations
- class: silence
- claim: C1 — "ignoring it can only over-admit, and an over-admitted cell surfaces as the bound refusal with the `guard.all` complement as the named remedy". Risks — "**`unless` on the stepped tag over-admits cells** → the bound refusal fires".
- evidence: the bound refusal is `conform`→`conformDomain` on the stepped literal (`internal/table/load.go:1869-1892`); it fires only when `cell + n` leaves `min..max` or the domain. An over-admitted INTERIOR cell — e.g. `attempt` declared `0..9`, rule `guard.all.attempt lt = 5` with `guard.unless.attempt eq = 3`, step 1 — yields a row for cell 3 writing 4, which conforms. No refusal fires.
- why it matters: the mitigation is stated as if it always catches the over-admission, and it catches only the boundary case. The interior over-admitted cell mints a row the author excluded, and its `guard.all eq = 3` atom now OVERLAPS whatever row the author wrote to claim cell 3 — `graph-overlap` at lint (blocking), which is a finding, not a refusal, and not the one C1 names. C1 should say the interior case surfaces as `graph-overlap`, or A5's "no overlap the hand-unrolled table did not have" must be scoped to rules with no `unless` on the stepped tag. A5 is `Pending` and its MVV fixture pair does not carry an `unless` on a stepped tag.

### F10 — the gate section is unlocked template text (mechanical, not design)
- anchor: 0030:§finalization-gate
- class: silence
- claim: the section carries the template's bracketed guidance verbatim, lines 1053-1172.
- evidence: `recs lint` in this evidence dir reports `gate:inline` plus five `placeholder:survived` advisories over 1083-1172; the `artifact` edge at line 1056 still points at the unexpanded token `{ARTIFACT_DIR}/gate.md` (the one edge in the record with `resolved=ABSENT`).
- why it matters: noted for completeness only — this is the pre-lock state the Finalization Gate itself resolves, not a design defect. It does mean the Gate's own "each cited `path::Symbol` resolves on `main`" check has not yet run, which is what would have caught F1.
