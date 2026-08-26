Model: claude-opus-5[1m]

# 3amigo iter-2 Persona 2 — Implementer (RDR cli/0011, delta-scoped)

Delta-scoped to the seven open entries carried from
`evidence/3amigo/persona-2-implementer.md`. Record read through
`rdr inspect`; source read at `internal/resolve/resolve.go`,
`internal/resolve/guard.go`, `internal/cli/flow_next.go`,
`internal/cli/flow_input.go`, `internal/cli/flow_exec.go`.

---

## 1. `uncomparable` fixture reachability — CLOSED

**Anchors:** `0011:C3` (third paragraph), `0011:S3`, `0011:MVV` step 6,
`0011:A3`, `0011:§failure-modes` F3, `0011:§mini-checks` `oracle` row
"MVV 6 / S3 conflicted".

Every site now names the two-reader owned-key collision and explicitly
refuses the `--tag` path:

- `0011:C3`: "The conflicted fixture MUST be built from two invoked readers
  writing the same owned key with differing values; a repeated `--tag`
  cannot reach the kernel (`flow_input.go::parseTags` refuses it as
  `flow-tag-duplicate`) … MUST NOT be accepted as this fixture."
- `0011:S3`, `0011:MVV` 6, `0011:§failure-modes` F3, and the `oracle`
  mini-check row all repeat the correction. No site still claims `--tag`.
- `0011:A3` even self-corrects the earlier `kernelTags` wording by name.

Producer verified reachable in source. `internal/cli/flow_exec.go:223`
inside `flowRequest.runReaders`:

```go
owned = append(owned, result.OwnedSnapshot()...)
```

— one append per invoked reader into a single `owned` slice, no dedup, only
a terminal `slices.SortFunc` by key (flow_exec.go:226). That slice becomes
`Input.Owned`, and `internal/resolve/resolve.go::merge` (resolve.go:235–243)
sets `conflicted` exactly when a key repeats **within one provenance** with
a differing value — which two readers' owned snapshots satisfy. The refusal
that blocks the `--tag` path is confirmed at `internal/cli/flow_input.go:256`
(`case seen[key]: … codeTagDuplicate`), upstream of `kernelTags`.

One residual, not a defect: nothing in the record states that the fixture
needs TWO read accessors bound to (at least) one artifact role each and that
`runReaders`' pre-flight (flow_exec.go:181–191) refuses `flow-artifact-missing`
for any invoked reader's unbound role. That is fixture mechanics, not a
contract gap, and `0011:S3` naming `flow_exec.go::runReaders` is enough to
find it.

---

## 2. Named kernel→CLI carrier for indeterminate match atoms — STILL-OPEN

**Anchors:** `0011:C1` (fourth paragraph, the CARRIER clause), `0011:C1`
(fifth paragraph, the `unknown` obligation), `0011:A1`,
`0011:§implementation-plan` Phase 2, `0011:§mini-checks` `authority` row
"candidate `unknown` list", `0011:§mini-checks` `disposition` table.

The carrier IS now named and it is coherent for the plain case. C1 fourth
paragraph: "the probe's disposition for such a row is that refusal, and its
`Refusal.Undecided` payload is the CARRIER … no new field on `Plan` is
required", plus the closing sentence "A row whose probe returns a `Plan` has
a fully decided match by construction, so it contributes no match facts."
That last sentence is the exact invariant my earlier finding asked for, and
it holds for the single-fault case.

It does NOT hold on the compound case, and `0011:A1` is what proves the gap
rather than closing it. A1 enumerates the probe's reachable dispositions as
"a plan, `no_match`, `owned_state_unavailable`, or `guard_unevaluable`" and
states the ordering as fact: "`internal/resolve/resolve.go::gate` returns
`owned_state_unavailable` then `guard_unevaluable` before that count."
Source confirms the precedence — `gate` (resolve.go:517–566) runs
`missingOwned` and returns `KindOwnedStateUnavailable` at resolve.go:531–537
**before** the `GuardUnevaluable` collection loop at resolve.go:540–566.

So take a row that is BOTH (a) indeterminate on a match atom and (b)
requires an owned key no invoked reader established. C1's fifth paragraph
obliges BOTH facts into `unknown`:

> "Every `indeterminate` match atom, every guard fact the kernel could not
> decide, every owned key no invoked reader established … MUST appear in
> that candidate's `unknown` list"

and its next-to-last sentence explicitly keeps the row a candidate:
"`guard_unevaluable` and `owned_state_unavailable` MUST leave the row a
candidate with those facts reported." But the probe returns
`owned_state_unavailable`, whose `Refusal` carries `MissingOwned []string`
(resolve.go:367–369) and an EMPTY `Undecided` — the guard_unevaluable branch
never ran. The match facts are not on the wire. The CLI cannot recompute
them: `TagSet.conflicting` is unexported (resolve.go:154) and `assembledView`
has no conflicted concept (`0011:A3`).

This is not the `Plan` hole from iter-1 — that one is genuinely closed — it
is the same hole moved one disposition over, and A1's own ordering sentence
is the evidence. The `disposition` mini-check table lists "match key
present, conflicted" and "owned key no reader established" as separate
independent rows and never crosses them; no row of that table, no MVV step,
and no scenario covers the conjunction.

Note the mirror question for the kernel selection site is also unstated:
where does C1's match check sit relative to `missingOwned`? If the
`indeterminate` exclusion happens in `Resolve`'s candidate loop (as Phase 1
says: "`Resolve`'s candidate loop and `escapeOrRefuse`'s escape-edge filter
are both SELECTION sites"), then an indeterminate-match row never reaches
`gate`, and its `RequiresOwned` obligation is dropped along with it — which
is precisely the pruning-masks-owned-state hazard `guard.go`'s block-boundary
comment (guard.go:129–135, citing `0007:C11` and D8) was written against.
The record does not say whether the match refusal preempts
`owned_state_unavailable` or the reverse, and the two orderings give
different refusal kinds on the same input.

**Blocks:** the Phase 1 site of the match check (inside `Resolve`'s loop vs.
folded into `gate` alongside the guard verdict), which fixes the refusal
precedence for the compound row; and Phase 2's `summarize`, which must
either read two payload fields or be told that this conjunction reports only
half its facts. Also blocks writing the compound-case oracle, which no
scenario currently owns.

---

## 3. Ordinary-selection-site blast radius — CLOSED

**Anchors:** `0011:C1` (second and third paragraphs),
`0011:D-verdict-use-by-site-class`, `0011:§consequences` (bullets 7 and 8),
`0011:S5`, `0011:A7`.

C1's second paragraph now makes the SELECTION/REPORTING split normative by
name ("a REPORTING site (`flow next`) MUST NOT exclude on it … a SELECTION
site (`flow resolve`, which commits a transition and its `Writes`) MUST NOT
select on it"), and the third paragraph states the disposition for both
kernel sites.

Both of my cases are answered, and answered the way the concern pushed:

(a) **Unescaped write-performing plan.** C1 third paragraph: "an ordinary
row whose match verdict is `indeterminate` MUST NOT be selected as a plan
… It instead refuses `guard_unevaluable`."
`0011:D-verdict-use-by-site-class` names the rejected alternative in exactly
those terms ("that lets `flow resolve` emit an unescaped write-performing
plan whose match precondition was never checked").

(b) **`ambiguous_match` against a genuine match.** C1 third paragraph: "and
MUST NOT count toward the `len(selected)` that yields `ambiguous_match`
against a row that genuinely matched." `0011:§consequences` carries the
consequence of that choice as a negative bullet ("a table whose rows overlap
only under undecidability resolves where a naive reading might expect a
refusal"), so the trade is booked rather than silent.

`0011:A7` remains `Pending` with Method "Source Search + MVV Test", but its
Evidence now says the four named oracles are "the claim's floor, not its
verified extent" and names both selection sites plus the five files to
enumerate — an honest Stage-6 obligation, not a silent gap. Source confirms
the two sites: resolve.go:472 (`if !view.matches(row.Match) { continue }` in
`Resolve`'s candidate loop) and resolve.go:614 (the same call in
`escapeOrRefuse`).

---

## 4. `Refusal.Undecided` vs. guard-path invariants (A8 scoping) — CLOSED

**Anchors:** `0011:A8`, `0011:C1` (third paragraph, the parenthetical
widening clause), `0011:§implementation-plan` Phase 1.

A8 is correctly scoped. It names both source sites and both are confirmed:
`guard.go:121` (`if !isGuardBlock(atom.Block) { … continue }`, whose own
comment says the skip is deliberate — "an atom that contributes no operand
cannot have blocked it", guard.go:137–139) and `guard.go:144` /
the `verdict != GuardUnevaluable` payload discard. `isGuardBlock` is
`b == BlockAll || b == BlockUnless` (guard.go:31–33), so `BlockMatch`
(guard.go:24) is excluded today exactly as A8 states.

The scoping is right in the load-bearing respect: A8's claim is not "match
atoms can ride `Undecided`" (that is a design choice C1 makes) but the
narrower, falsifiable "while leaving every existing guard payload
byte-identical", and its Evidence names the specific hazard — "`0007:C8`
pins per-atom completeness and the payload sort (`compareUndecidedAtoms`
orders by block), so a new block value participates in ordering". The
`If wrong` arm routes correctly (separate `Refusal` field; C1's carrier
clause and `D-undecided-reporting-shape` change with it). C1 books the
widening explicitly rather than leaving it implied, and disclaims the
`RefusalKind` set change (`0007:REQ-79`).

One thing A8 does not cover, which I flag only because it is inside A8's own
"byte-identical" claim: a row that carries BOTH a decidable-false guard and
an indeterminate match. `evaluateAtoms` returns `GuardFalse` and `gate`
prunes it (resolve.go:522–524) — correct under C1 (guard-decided-false
excludes in both modes), so no payload is owed. Not a gap; noting it so the
Stage-6 verifier does not re-derive it.

---

## 5. `unknown` dedup rule and the `not-evaluated` token — CLOSED

**Anchors:** `0011:C1` (fifth paragraph), `0011:D-undecided-vocabulary`,
`0011:D-undecided-reporting-shape`, `0011:§mini-checks` `authority` row.

Dedup is now stated and workable. C1: "Entries are deduplicated by
`{key, reason}` pair, not by key, so a key that is both an undecided match
atom and an unestablished owned key reports once per distinct reason." That
answers my earlier question directly (one entry or two? — two, when the
reasons differ; one when they coincide) and is implementable against the
shipped `slices.Contains(c.Unresolved, …)` dedup at flow_next.go:180/190/199
by swapping the element type. The `authority` mini-check row repeats the same
rule ("Dedup key is the `{key, reason}` pair, not the key"), and
`D-undecided-reporting-shape` pins the ordering ("sorted by `(key, reason)`")
so the merge of kernel-supplied and CLI-computed facts is stable.

`not-evaluated` is coherent. C1 argues it from the closed set rather than
asserting it: a un-run gate id "is neither `absent` (the gate is declared and
its id is reported) nor `uncomparable` (a gate id is not a view key with a
value)". The scoping is explicit and testable — "which this RDR mints for the
gate class ONLY and which never appears on a match or owned-key entry" —
and the seam/reporting-surface split is stated ("`0007:C8`'s two-member set
is unchanged for the seam it governs; `unknown` is a CLI reporting surface
that spans three sources"). `0011:C3` even converts it into a build rule:
"a re-homing that has to invent a reason token is not mechanical and is a
defect."

Caveat inherited from finding 2, not re-raised here: the *rule* is fine; what
is unavailable on the compound row is one of the *inputs* to it.

---

## 6. `--all` probe construction and the `summarize` BlockMatch filter — CLOSED

**Anchors:** `0011:C2` (first paragraph), `0011:§implementation-plan`
Phase 3, `0011:S2`, `0011:A2`.

Both halves settled, and settled against the honest reading. Construction:
C2 says the exclusion-free predicate is "achieved by the probe omitting the
match pattern so the kernel never sees those atoms" — Phase 3's "route it to
a probe with the match pattern omitted", i.e. today's `probe.Match = nil`
(flow_next.go:278).

The `summarize` question is answered explicitly and reclassified: C2 now
states "This is a DEPARTURE from shipped behaviour, not a preservation of
it", cites A2 for why (`summarize` walks `Row.Atoms`, which carries
match-block atoms — confirmed at flow_next.go:177–184, a loop with no
`Block` test), mandates the filter ("the `--all` branch MUST filter
BlockMatch atoms out of `unknown`"), mandates an oracle ("an oracle MUST
assert their absence"), and fences it off from A4 ("it lies outside A4's
predicate-scoped classification"). `0011:S2`'s Expected carries the paired
assertion (`--all` reports no match entry; default reports `{key, absent}`).

---

## 7. `--all` rejection mechanism for resolve/read-state/set-state — CLOSED

**Anchors:** `0011:C2` (final paragraph), `0011:S2`,
`0011:§mini-checks` `trace` step 1.

The mechanism is named: "satisfied by NON-REGISTRATION (registered on next
only, not on the shared `registerSelectionFlags`), so cobra's shipped
unknown-flag path supplies the error and exit code. This contract mints no
new refusal code for it." The `trace` mini-check step 1 names the converter
site (`internal/cli/root.go`), and `0011:S2` states the oracle shape:
"asserted on that error class and exit code rather than a `flow-*` code."

Source is consistent: `newFlowNextCmd` calls the shared
`registerSelectionFlags(cmd)` (flow_next.go:78) and registers
`--evaluate-gates` on itself (flow_next.go:80–81), so `--all` follows the
`--evaluate-gates` precedent exactly. The iter-1 source-anchor ambiguity is
resolved by C2 naming `registerSelectionFlags` as the site NOT to use.

---

## Delta summary

CLOSED: 1, 3, 4, 5, 6, 7. STILL-OPEN: 2 (compound indeterminate-match +
missing-owned row: the carrier is preempted by `owned_state_unavailable`,
and the match-check-vs-`missingOwned` precedence at the kernel selection
site is unstated). No NEW-GAP opened by the rewrite.
