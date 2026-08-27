Model: claude-opus-5[1m]

# 3amigo — resolve, RDR 0010

Origin ledger = the 20 findings in [`consolidation.md`](consolidation.md).
Every finding exits exactly one way. Grounding gate run against code on `main`,
`{RDR_RESOURCES}`, and the RDR's own decided text before any edit.

## Dispositions

| Finding | Disposition | Origin | Section touched |
| --- | --- | --- | --- |
| P1-1 | **fixed** | P1-1 (hotspot `C5`, 3 personas) | `C5` — zero-participating-dimension groups take `graph-unprovable-coverage`; Risk entry re-based from documentation to contract |
| P2-1 | **fixed** | P2-1 (hotspot `C4`×`S5`) | `C4` — text mode de-specified to 0005's generic renderer; no per-verb special case |
| P2-2 | **fixed** | P2-2 | `C1` — agreement check moved to a step at/after `loadTags`; refusal precedence stated |
| P2-3 | **fixed** | P2-3 | `C1` — check made one-directional; `§consequences` census added |
| P3-1 | **fixed** | P3-1 (same passage as P2-1) | `C4`, `S5` — assertion pinned to the renderer as built (`emit.<key>: <value>`, `emit: (none)`) |
| P1-2 | **fixed** | P1-2 | `§phase-4-model-and-docs` — authoring-doc Done criterion named |
| P2-4 | **fixed** | P2-4 | `C3` — `Emit []EmitValue` (new type), not `TagValue`; `expand` carry stated |
| P2-5 | **fixed** | P2-5 | `C3`, `S3` — dedup clause cut as unreachable; duplicate-key test case removed |
| P2-6 | **fixed** | P2-6 | `C3`, `S3` — `map[string]string` + decoder type arm ⇒ `malformed TOML`, category-only assertion |
| P2-7 | **fixed** | P2-7 (hotspot `C5`) | `C5`, `§technical-design`, `§authority` — predicate augments `len(Initial)`; accessor with `state-machine` zero value |
| P2-8 | **fixed** | P2-8 | `C2` — real refusal route named (`checkAccessorBindings`); the apparent gap closed at `loadAccessors` |
| P2-9 | **fixed** | P2-9 | `C3` — `emit` appended last; raw-string cell, not `renderValue` |
| P3-2 | **fixed** | P3-2 | `S6` — old-binary arm removed from the standing suite; A8 stays a discharged spike |
| P3-3 | **fixed** | P3-3 | `§testing-strategy` Done clause — licensed-diff rule replaces the adjective |
| P3-4 | **fixed** | P3-4 (hotspot `MVV`) | `S4`, Oracle row 1 — match-only control now asserts a positive finding |
| P1-3 | **fixed** | P1-3 | `MVV` — consumer-side acceptance + tracker split (intrastate#zdat vs rdr#tmxk) |
| P2-10 | **fixed** | P2-10 | `C3`, `§technical-design` — false per-row uniqueness claim replaced with the rule-scoped join argument |
| P2-11 | **dismissed-with-cite** | P2-11 | Re-raise of `0010:BR4` (`emit` in the `flow next` candidate preview, already rejected). The decision was implicit in C4's silence; made explicit in `C4` rather than re-litigated. |
| P3-5 | **fixed** | P3-5 (hotspot `C5`) | `C5` — class-silence bounded to an assertable form (no `Code`/`Element`/`Message`/`Detail` naming the class) |
| P3-6 | **fixed** | P3-6 | `C1`, `S1` — detail token fixed as `owned=<n>` |

**19 fixed · 1 dismissed-with-cite · 0 charted · 0 needs-tiebreaker.**

## Grounding-gate corrections (findings whose cited mechanism was checked, not assumed)

Two edits were corrected mid-pass after reading the source rather than
accepting a finding's framing:

- **P2-8** — a first draft asserted the `[initial]`-plus-writer gap closes via
  C1 because "the writer makes the tag owned". False: `Provenance` is a
  declared `[tags.<tag>]` field, never derived. The real closure is earlier and
  harder — `internal/table/load.go::loadAccessors` refuses any writer naming a
  non-owned tag (`CatWriteToNonOwnedTag`, `0002:C2`). Corrected in place.
- **P2-6** — a first draft said a non-string emit value refuses
  `unknown schema field`. False: `decodeStrict` maps only
  `toml.StrictMissingError` (an unknown *key*) to `CatUnknownSchemaField`; a
  well-named key of the wrong type falls to `CatMalformedTOML`. Corrected in
  `C3` and `S3`.

## Compute-don't-argue

- `§consequences`' "every existing model loads unchanged" is now **measured**:
  a census over checked-in `[model]`-carrying TOML under `internal/` and
  `models/` returns **zero** models with an empty owned set, so C1's
  decision-table arm cannot fire on any of them and (one-directional) the
  state-machine arm has no refusal to fire at all.
- Cross-lens check: cove pinned `C5`'s terminal-arm scoping in the prior pass;
  this pass amends `C5` again (zero-dimension fence, root predicate). The two
  edits are compatible — cove scoped the `0006:C18` override to the root arm,
  this pass augments the seeding predicate and adds an invariant-4 arm. Neither
  re-pins a field the other set.

## Needs (re)verification — carried to Stage 6

- **A13 (new, Pending)** — a decision-table group with zero participating guard
  dimensions can take `graph-unprovable-coverage` without a taxonomy change and
  without disturbing the state-machine class. Method: MVV Test. Booked because
  C5's new fence is a load-bearing claim about an emission that has not been
  executed: `emitStructurallyUnprovable` is per-dimension today and would emit
  with an empty `Dimension` field, and no fixture sweep has confirmed that no
  checked-in state-machine model has a zero-dimension group.
- **A10 (Verified, unchanged in substance)** — its code partition still holds;
  the claim gained a sentence noting `graph-unprovable-coverage` keeps its code
  but gains a class-conditioned arm (pointing at A13). No flip to Pending: the
  fourteen-code partition A10 asserts is untouched.
- **A12 (Pending, unchanged)** — carried from cove.
- No previously Verified assumption was invalidated.

## Charted to successor

None. Every finding landed inside this RDR's existing contract surface. The
one candidate for net-new scope — P1-1's fence — was resolved *inside* `C5`
using an existing taxonomy code (invariant 4's `graph-unprovable-coverage`)
rather than by adding a contract or a code, so it is an amendment to a decided
contract, not an expansion. `0010:BR6`'s rejection of a class-announcing
finding was checked and does not cover it: BR6 rejects announcing the *class*,
this announces an unprovable *group*.

## Convergence

One iteration. The lens ran once, all 20 findings dispositioned, no finding
re-opened by another finding's fix. The amendment sweep (§amendment-sweep) ran
per fix: the `key=value` text-mode change swept 5 candidate sites (2 text-mode,
3 dump — the dump sites are a different claim and stayed), the vacuous-coverage
fence swept 6 sites, the one-directional C1 change swept the class-omitted
control in `MVV` step 6, `S4`, and Oracle row 2, and the class-keyed site count
went 3 → 4 in both the Authority cue and the drift Risk. No re-run owed: the
fixes are pins and de-specifications against a stable frame, not a rewrite.
