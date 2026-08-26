Model: claude-opus-5[1m]

# 3amigo Persona 2 — Implementer (RDR cli/0011)

Owned set: `0011:C1`, `0011:C2`, `0011:C3`, the six `0011:D-*` decisions, and the
`source-anchor` edges. I widened past that set three times, and say so at each
finding: to `0011:MVV` / `0011:S3` (a contract obligation whose *test* is
unreachable — the silence is in C1, but only MVV names the input), to
`0011:A3`/`0011:A7`/`0011:§consequences` (C1's second paragraph asserts a
kernel-wide seam change whose blast radius is only described outside the
contract), and to `0011:§implementation-plan` Phase 2 (C1 states an obligation
with no named carrier; the Phase intent is where the carrier would live).

Read against source at `internal/resolve/resolve.go`, `internal/resolve/guard.go`,
`internal/cli/flow_next.go`, `internal/cli/flow_input.go`,
`internal/cli/flow_exec.go`.

---

## 1. BLOCKING — The `uncomparable` fixture C3 mandates is unreachable through the CLI

**Anchors:** `0011:C3` (third paragraph, the discriminating-fixture obligation),
`0011:C1` (the `uncomparable` reason), `0011:MVV` step 6, `0011:S3`,
`0011:§mini-checks` (`disposition` row "match key present, conflicted";
`oracle` row "MVV 6 / S3 conflicted"). Widened to MVV/S3 because C3's obligation
names the fixture but only MVV names the input that produces it.

C3 requires the build add "a row whose match key is supplied twice with
differing values (candidate in both, reason `uncomparable`)". MVV 6 and the
`disposition` table both realise that input as a repeated supply. But
`internal/cli/flow_input.go::parseTags` **hard-refuses a repeated `--tag` key
before any accessor runs**:

```go
case seen[key]:
    return nil, userErr(codeTagDuplicate, key,
        "the tag `"+key+"` is given more than once")
```

So `flow next --tag k=x --tag k=y` never reaches `kernelTags`, never reaches
`assemble`, and can never produce `tv.conflicted`. `0011:A3`'s evidence sentence
— "`internal/cli/flow_exec.go::kernelTags` passes a repeated `--tag` through
without dedup" — is true of `kernelTags` in isolation and false of the path: the
refusal is one layer upstream, in `parseTags`.

Reading source for the reachable path: `TagSet.merge` marks conflicted only for
a repeat **within one provenance**. Observed comes solely from `--tag` (blocked).
That leaves **owned** — `flowRequest.runReaders` appends
`result.OwnedSnapshot()...` across *multiple readers* into one `owned` slice with
no dedup, so two readers reporting the same key with different values is the one
CLI-reachable conflicted key. But an owned key cannot be a `--tag`, and C3/MVV 6
describe the fixture as a supplied observed tag.

**Clarification requested:** which is it —
(a) does `parseTags` relax to admit a repeated `--tag` (a change to a shipped
0005 refusal `flow-tag-duplicate`, which C1/C2/C3 nowhere authorise and which
`0011:§scope-verification` does not scope), or
(b) is the CLI-level `uncomparable` fixture a *two-reader owned-key collision*
(a fixture shape nothing in the record describes and which requires a
multi-reader model), or
(c) is the `uncomparable` arm exercised only at the `internal/resolve` package
level (a `TagSet` built directly), with the CLI-level MVV 6 conflicted row
struck?

**Blocks:** writing `0011:S3`'s third fixture and `0011:MVV` step 6 at all —
the entire `uncomparable` half of C1's closed reason set has no reachable
`flow next` input. Also blocks the `disposition` mini-check row and the `oracle`
row that names it "the case the pre-Stage-4 design got wrong": that oracle
cannot be built as written, so C1's most load-bearing distinction ships
untested at the CLI seam.

---

## 2. BLOCKING — C1 obliges the CLI to report indeterminate match atoms, but the probe returns a `Plan`, which carries no atom payload

**Anchors:** `0011:C1` (third paragraph: "Every `indeterminate` match atom …
MUST appear in that candidate's `unknown` list"), `0011:§implementation-plan`
Phase 2 ("carry the kernel's undecided match atoms into the candidate's
`unknown`"), `0011:§mini-checks` `authority` row "candidate `unknown` list …
folds three sources: kernel indeterminate atoms". Widened to Phase 2 and the
mini-check because C1 states the obligation and neither names the carrier.

Trace the probe as C1 and Phase 2 specify it (`excluded` stops stripping
`probe.Match`), for the case C1 exists to fix — a row whose match key is
ABSENT and whose guard is decidable-true:

1. `Resolve` candidate loop: match verdict is `indeterminate`; C1 says it MUST
   NOT exclude, so the row survives.
2. `gate` decides the guard true, no owned gap → `selected` has 1 row.
3. `switch len(selected) { case 1: return Result{Plan: planOf(...)} }`.

The disposition is a **`Plan`**. `resolve.Plan` (resolve.go:338–351) has exactly
six fields — `RuleID`, `SourceLocator`, `NextTags`, `Writes`, `Revision`,
`Escaped`. There is **no `Undecided` field and no atom payload**. `Refusal` has
`Undecided []UndecidedRow`, but the probe did not refuse. `TagSet.conflicting`
is *unexported* and `Lookup` gives presence but not conflictedness, so the CLI
cannot recompute the verdict either (this is exactly `0011:A3`'s point, and
`0011:D-placement` rejects the CLI recomputation on the same ground).

So on the happy path — the one MVV 6 "absent (candidate, `{key, absent}`)"
asserts — the CLI has **no channel** to learn which atoms were indeterminate.

**Clarification requested:** what is the new kernel surface? Concretely, one of:
(a) a field on `Plan` (e.g. `Plan.Undecided []UndecidedRow`) — an addition to a
success disposition that `0009:C3`'s "exactly one plan or exactly one refusal"
prose does not obviously admit, and which `flow resolve` would then have to
decide whether to surface or drop;
(b) an exported three-valued entry point such as
`func (s TagSet) MatchVerdict(want []Tag) (Verdict, []UndecidedAtom)` plus an
exported `TagSet` constructor — but `TagSet` is not constructible from outside
`internal/resolve` today (`assemble` is unexported), and `0011:D-placement`
explicitly rejects "a separate `resolve.Candidates` entry point";
(c) something else.

Note this also fixes the shape of `0011:D-undecided-reporting-shape`: the CLI's
`{key, reason}` pair must be derivable from whatever the kernel hands back, and
`UndecidedAtom` carries `Key`, `Block`, plus operator/literal — so the CLI is
either projecting `UndecidedAtom` down to two fields or the kernel returns
something narrower.

**Blocks:** the entire Phase 1 API design and therefore Phase 2. I cannot write
`internal/resolve`'s signature on Monday, cannot write the Phase 1 three-valued
oracles that "replace `TestReq78`", and cannot write MVV 6's `{key, absent}`
assertion.

---

## 3. BLOCKING — C1's kernel clause silently changes ordinary `flow resolve` selection, and only the escape half is accounted for

**Anchors:** `0011:C1` (second paragraph, "The verdict is the SEAM's, so it MUST
apply at every site that asks it"), `0011:D-placement`, `0011:A7`,
`0011:§consequences`. Widened to A7 and Consequences because C1's paragraph
asserts uniformity without stating the ordinary-site consequence, and A7 and
Consequences are where the blast radius is described — describing only half of it.

C1 makes `indeterminate` non-excluding at **both** `resolve.go:472` (ordinary
candidate selection) and `resolve.go:614` (escape rescue). `0011:A7`'s statement
is scoped entirely to `escapeOrRefuse`; `0011:§consequences`'s negative bullet
likewise says "the escape path widens too". Nothing in the record states the
ordinary-site consequence, which is strictly larger. Two concrete cases from
shipped source:

**(a) An unescaped plan on a row whose match key was never established.**
`guard_atoms_test.go::TestReq78_MatchPatternStillFoldsAbsenceIntoNonMatch`
(guard_atoms_test.go:930) builds `guardedRow` + `Match{absentKey: "3"}` plus a
`no_match` escape row and asserts the disposition is the **escaped** plan. Under
C1 the ordinary row is no longer excluded, its guard passes,
`len(selected) == 1`, and `flow resolve` emits an **unescaped plan whose match
precondition was never checked** — it will `Writes` owned state on a transition
selected under an admittedly-undecidable match. The record names retiring
`TestReq78`, but frames it as closing `0007:REQ-78`; it does not say that the
replacement disposition is a *write-performing unescaped plan* rather than a
rescue.

**(b) A previously-clean plan becomes `ambiguous_match`.** Where a table has one
row that truly matches and a sibling row for the same outcome with an
indeterminate match atom, both now survive to `gate`; if both pass,
`len(selected) == 2` → `escapeOrRefuse(KindAmbiguousMatch, ...)`. A resolution
that used to succeed now refuses. `0011:A1` proves `ambiguous_match` is
unreachable **on a one-row probe** and calls the sub-question "vacuous" — true
for `flow next`, and exactly not true for `flow resolve` over a real table.

**Clarification requested:** is (a) intended — may an ordinary (non-escape)
plan be selected while one of its own match atoms is `indeterminate`, and does
it then perform its `Writes`? If yes, C1 should say so and Consequences should
carry it, because it is the sharpest behaviour change in the RDR. If no, C1
needs a clause distinguishing the *selection* sites (where `indeterminate` may
have to disqualify, or raise a typed refusal) from the *reporting* site
(`flow next`, where it must not). And for (b): should an indeterminate-match row
be allowed to induce `ambiguous_match` against a row that genuinely matched?

**Blocks:** Phase 1's one-line description — "`Resolve`'s candidate loop and
`escapeOrRefuse`'s escape-edge filter both exclude only on `no-match`" — is the
whole of the guidance, and it does not tell me what disposition to return in the
two cases above. It also blocks classifying the `internal/resolve` suite
(`0011:A7` is still `Pending`, and its Method scopes the enumeration to *escape*
oracles only — the ordinary-selection oracles are not in its census).

---

## 4. HIGH — `Refusal.Undecided` is documented as a `guard_unevaluable`-only payload; C1 makes match atoms ride it

**Anchors:** `0011:C1` (second paragraph, "its undecided atoms ride the
refusal's undecidable-row facts, the same way an ordinary row's do"),
`0011:D-undecided-vocabulary`.

`resolve.Refusal.Undecided` is field-documented as populated **"on
`guard_unevaluable`"**, and `evaluateAtoms` explicitly discards payload entries
for any row whose verdict is not `GuardUnevaluable`:

```go
if verdict != GuardUnevaluable {
    // Only an undecidable row reports a payload (`0007:C8`); a decided
    // row's entries are discarded.
    return verdict, nil
}
```

C1 requires match atoms to ride "the refusal's undecidable-row facts" for an
indeterminate escape edge. But an escape edge with an indeterminate match and a
**decided** guard produces no `UndecidedRow` at all today, and the refusal it
would ride (per C1, a `no_match` when nothing rescues) is a kind on which
`Undecided` is currently always empty.

**Clarification requested:** does `Refusal.Undecided` become populated on
`no_match` (and possibly `ambiguous_match`) too? If so, `UndecidedAtom.Block`
would carry `BlockMatch` — which `guard.go::isGuardBlock` currently and
deliberately excludes from the payload ("an atom that contributes no operand
cannot have blocked it"). Does `0007:C8`'s payload contract widen, or does C1
mint a separate `Refusal` field for match facts?

**Blocks:** the `Refusal` struct shape, and `0011:S5` ("`flow resolve` over the
same seam change") — I cannot assert on a payload whose field I do not know.

---

## 5. HIGH — C1's `unknown` obligation for owned keys collides with the probe returning `owned_state_unavailable`

**Anchors:** `0011:C1` (third paragraph: "every owned key no invoked reader
established … MUST appear in that candidate's `unknown` list"; and
"`owned_state_unavailable` MUST leave the row a candidate with those facts
reported"), `0011:§mini-checks` `authority` row for the `unknown` list.

C1 says the `unknown` list folds three sources, and the `authority` mini-check
agrees: "kernel indeterminate atoms, absent owned keys, un-run gate ids". Two of
those three the CLI computes itself in `summarize` from `row.RequiresOwned` and
`row.Gate` against `assembledView` — a **flat string map with no conflicted
concept** (`0011:A3` establishes this: zero occurrences of `conflicted` in
`internal/cli`). The third now comes from the kernel probe.

So `unknown` will mix a CLI-derived `absent` (from `assembledView` presence)
with a kernel-derived `absent`/`uncomparable` (from `TagSet`). These two views
disagree by construction: `assembledView` reports a conflicted-owned key as
*present* (it just overwrote), while the kernel reports it `uncomparable`.

**Clarification requested:** when a row's `RequiresOwned` key and its match atom
name the **same key**, and the key is present-but-conflicted, does the candidate
get one entry or two, and with which reason? And is the owned/gate half of
`unknown` re-sourced from the kernel too (making `assembledView` reporting-only
for `observed`/`owned` echo), or does it stay CLI-computed with the known
divergence?

**Blocks:** `summarize`'s implementation and the dedup rule — the shipped code
dedups by `slices.Contains(c.Unresolved, key)` on a bare string; with
`{key, reason}` pairs the dedup key is now ambiguous (key alone, or the pair?).
C1 does not say, and MVV 7's "`unknown` present as `[]`" assertion does not
reach it.

---

## 6. MEDIUM — `--all` is specified as a predicate, not as a probe construction, and the two readings differ

**Anchors:** `0011:C2` ("the candidate predicate MUST be exactly the one
`0005:C1` specified … with the row's match pattern taking no part"),
`0011:§implementation-plan` Phase 3 ("route it to a probe with the match pattern
omitted"), `0011:D-selection-predicate`.

Phase 3's "probe with the match pattern omitted" is today's `probe.Match = nil`,
and it is the natural reading. But C2 also says "A match atom is then neither an
exclusion nor an entry in `unknown`, **whatever its key's presence**". Under the
`probe.Match = nil` construction that falls out for free — the kernel never sees
the atoms. Under the alternative construction (retain match, ignore the verdict)
the CLI would have to actively suppress them.

More consequentially: `summarize` currently walks **`row.Atoms`**, which
`internal/table/model.go::Row.KernelRow` routes by `Block` — so it carries
match-block atoms alongside guard atoms (this is `0011:A2`'s finding). Today
match-block atoms over an absent key **already land in `unresolved`**. So under
`--all`, C2's "neither an exclusion nor an entry in `unknown`" is a *removal* of
shipped 0005 behaviour, not a preservation of it — yet C2's framing is
"everything else `0005:C1` requires … MUST hold identically" and `0011:A4`'s
BREAKS-0 result is predicate-scoped.

**Clarification requested:** under `--all`, must `summarize` now **filter out**
`BlockMatch` atoms from `unknown`? If yes, that is a shipped-behaviour change in
the `--all` mode that C2 presents as 0005-preserving, and it needs an oracle. If
no, C2's "neither … an entry in `unknown`, whatever its key's presence" is
contradicted by `row.Atoms`.

**Blocks:** the `--all` branch of `summarize`, and whether `0011:S2`'s `--all`
half asserts a non-empty or match-free `unknown`.

---

## 7. MEDIUM — `--all` flag placement and the "must not be accepted by" clause have no stated mechanism

**Anchors:** `0011:C2` (final sentence: "`--all` MUST NOT be accepted by
`flow resolve`, `flow read-state`, or `flow set-state`"), `0011:D-naming`,
source-anchor edge to `internal/cli/flow.go::registerSelectionFlags`.

`registerSelectionFlags` is the shared per-verb flag registrar the source-anchor
edge names. Registering `--all` there would give it to every verb; registering
it in `newFlowNextCmd` gives it to `next` only. Phase 3 says "register `--all`
on `next` only", which reads as the latter — in which case the source-anchor to
`registerSelectionFlags` is misleading, and the "MUST NOT be accepted" clause is
satisfied by absence rather than by an explicit rejection.

**Clarification requested:** is the obligation satisfied by non-registration
(cobra then errors `unknown flag: --all`, an exit code the record does not
name), or does C2 want an explicit typed refusal with a `flow-*` code? These
produce different exit codes and different oracles.

**Blocks:** writing the negative oracle for C2's last sentence — `0011:S2` does
not cover it, and the `trace` mini-check step 1 asserts it under "C2 (`--all`
rejected by `resolve`/`read-state`/`set-state`)" with witness "flag name free:
no `"all"` in `internal/cli`", which is evidence of availability, not of
rejection semantics.

---

## 8. LOW — Phase 4's six mechanical re-homings are counted but not located precisely enough to be mechanical

**Anchors:** `0011:C3` (second paragraph), `0011:§implementation-plan` Phase 4.

C3 names the count and the files (`flow_next_0005_test.go` ×4,
`flow_adversarial_0005_test.go` ×1, `flow_mvv_0005_test.go` ×1) but not the
struct field. The shipped Go field is `candidate.Unresolved []string` with tag
`json:"unresolved"` (flow_next.go:49), and the file's own header comment at
flow_next.go:14 names `unresolved` as well. That is a seventh and eighth site
(field name + header comment) beyond the six *reads*, plus `summarize`'s three
`slices.Contains(c.Unresolved, …)` uses.

**Clarification requested:** is "six shipped reads" a count of *test
assertions* only (excluding the production field rename and the doc comment), or
does the audit claim six total? A build that renames the field will touch more
than six lines, and C3 says these "MUST NOT be counted against or excused by
A4's predicate-scoped BREAKS-0 result" — so the count is load-bearing for the
gate.

**Blocks:** nothing on its own; it blocks reporting Phase 4 complete against a
number I would otherwise contradict.

---

## 9. LOW — Finalization Gate sections are unauthored template text

**Anchors:** `0011:§contradiction-check`, `0011:§scope-verification`,
`0011:§cross-cutting-concerns`, `0011:§proportionality` (all four still carry
the bracketed template instructions rather than written responses).

Not an implementer clarification per se, but `0011:§scope-verification` is the
section that would tell me whether MVV 6 is in scope for the build or deferred —
which is directly upstream of finding 1. As written it tells me nothing.

**Blocks:** knowing whether the unreachable `uncomparable` fixture (finding 1)
is a scope deferral or a defect.
