Model: claude-opus-5[1m]

# 3amigo Persona 1 — Product Manager (RDR cli/0011)

Scope note: I read the owned set (`0011:§problem-statement`, `0011:§approach`,
`0011:§decision-rationale`, `0011:MVV`) first. Three of the four findings below
required widening, and I name what sent me at each. The recurring driver is that
the user outcome in `0011:§problem-statement` is stated for a *human/skill caller
reading a list*, while every acceptance criterion in `0011:MVV` and
`0011:§testing-strategy` is stated for a JSON payload — so the outcome's own
success condition has no line range inside the owned set.

---

## F1 (High) — The stated user outcome is "a skill can constrain its next step", but no criterion anywhere tests that a skill can

**Anchors**: `0011:§problem-statement`, `0011:MVV`, `0011:§consequences`,
`0011:§testing-strategy` (widened from MVV: the outcome's acceptance condition is
not in the owned set, which is what sent me).

`0011:§problem-statement` frames the outcome in the caller's terms: a skill author
"cannot constrain a skill's next step from it", and cites 0005's own materialized
premortem, "next omits enough condition detail to constrain a skill."
`0011:§consequences` restates it as "a skill can read the next step from the
candidate list at `stage=resolved` (A6) instead of the whole table."

But every step of `0011:MVV` and every scenario in `0011:§testing-strategy`
asserts on set cardinality and membership: 3 vs 21, presence/absence of a named
rule, `unknown` as `[]`. Cardinality is a *proxy* for the outcome, not the
outcome. Nothing in the record states what number of candidates constitutes
"constrained" for the actual consumer, nor asserts that the 3-row result at
`stage=resolved` is actionable — three candidates across three distinct outcomes
(`prelock`, `resolve-route-back`, `resolve-abandon`, per `0011:MVV` step 3) is
still an ambiguous next step for a skill that must pick one. `0011:§approach` is
explicit that the verb "never chooses among the rows that survive," and
`0011:D-selection-predicate` re-affirms it. So the shipped end state for the
motivating scenario is: the skill still cannot pick a next step from `flow next`;
it has a shorter menu.

That may be the right product decision — but the record never says so. It asserts
the outcome is delivered (`0011:MVV` End-state: "answers 'what is legal from
here' with the rows the supplied state can take") without ever reconciling it
against the un-narrowed remainder.

**Blocks**: the accept/reject decision at lock — a reviewer cannot tell whether
21→3 discharges the problem statement or merely improves it, because no
criterion states what "constrained enough" is. Prevents the only test that would
falsify the product claim: an end-to-end assertion that the motivating consumer
(rdr#tmxk / the skills, named in Related Issues and `0011:§background`) can now
select its next action, or an explicit statement in `§problem-statement` that
narrowing-not-selecting is the whole intended outcome and residual ambiguity is
`flow resolve`'s.

---

## F2 (High) — The primary reported symptom ("wall of candidates") survives the fix in the common case, and the record treats that as accepted rather than deciding it

**Anchors**: `0011:§problem-statement`, `0011:§decision-rationale` (premortem
paragraph), `0011:§risks-and-mitigations`, `0011:F2` (widened to Failure Modes
and Risks because the premortem in my owned `§decision-rationale` names the
symptom but routes its disposition outside the owned set).

`0011:§problem-statement` identifies the observable defect as "the candidate list
is the whole alphabet regardless of the supplied state." `0011:C1` then makes
`indeterminate` (key absent from the view) NOT exclude. `0011:F2` and the
`0011:§risks-and-mitigations` first bullet both concede the consequence directly:
"every row is a candidate with the same key under `unknown`, reason `absent`."
`0011:§decision-rationale`'s premortem names the same outcome — "saw every row
listed as a candidate with the key under `unknown` — the same wall-of-candidates
symptom as before, now with a hint" — and dismisses it as "the intended
three-valued behaviour."

From a product view this is the load-bearing gap: the record never establishes how
often the absent-key case is the *normal* case rather than the misconfiguration
case. `0011:A5` supplies direct evidence that it is normal for a whole model
class: over 0010's decision-table class with no `--tag`, "all 4 ordinary rows
reported, each carrying its observed match keys ... in `unresolved`" — i.e. zero
narrowing, full wall, by design. The A6 spike shows narrowing only because
`models/rdr.toml` matches on an *owned* key a reader supplies. So the fix delivers
the outcome exactly when the match keys are reader-owned, and delivers the old
symptom (plus a hint) when they are observed/tag-supplied. The record contains
both facts and never joins them into a statement of expected user experience.

The mitigation offered — "the fix is on the caller's side (bind the reader /
supply the tag) and is visible" — is a diagnosis, not the outcome the problem
statement promised.

**Blocks**: the go/no-go on whether this RDR fixes the reported defect for the
consumers that filed it. Prevents an acceptance test over the observed-key /
no-tag path stating what the user should see and whether that counts as success.
A one-clause statement in `§problem-statement` or `§consequences` — "the outcome
is delivered where match keys are owned; where they are observed and unsupplied
the caller must supply them, and the `unknown` reason tells them to" — would
close it.

---

## F3 (Medium) — The consumer migration is asserted as costless on a scoped result, and the caller-facing cost is never sized

**Anchors**: `0011:§consequences` (Negative: "callers scripted against the
enumeration must add `--all`"), `0011:MVV` step 8, `0011:C3`, Metadata
`Overrides`, `0011:A4` (widened from Consequences; the override cost is asserted
in prose in my owned sections and its evidence lives in A4/C3).

The record's breaking-change story is: A4 verified the 0005 suite's MOVES and
BREAKS classes are both empty, so the override "is a forward guard" (Metadata
`Overrides`). But A4 is explicitly *predicate-scoped* — `0011:C3` itself carves
out six mechanical `unresolved` → `unknown` re-homings and warns they "MUST NOT
be counted against or excused by A4's predicate-scoped BREAKS-0 result." I
confirmed the six against source: 4 in `internal/cli/flow_next_0005_test.go`
(lines 103, 163, 166, 404), 1 in `internal/cli/flow_adversarial_0005_test.go`
(446), 1 in `internal/cli/flow_mvv_0005_test.go` (66). C3's count is accurate.

The gap is that all of this is measured over *this repo's tests*. The candidate
payload is a published consumer contract — `internal/cli/flow_next.go:49` emits
`json:"unresolved"`, and `docs/cli-output-contract.md` documents `flow next
--as=json` as a scripting surface. This RDR simultaneously (a) flips the default
predicate and (b) renames and retypes a payload field, and no owned section
states the combined blast radius for an out-of-repo consumer. `0011:C2` requires
payload shape identity *between modes*, which does not help a consumer parsing
the old shape: `--all` restores the old candidate SET but not the old field NAME
or ELEMENT TYPE. So the documented escape hatch does not actually restore the old
contract, and the record nowhere says that.

**Blocks**: the release/versioning decision — whether this ships as a breaking
change with a migration note, and whether `--all` is honestly describable to users
as "0005's behaviour." Prevents a compatibility test asserting what a 0005-era
consumer sees under `--all`. Also note `0011:§cross-cutting-concerns` is unfilled
template while "versioning" and "incremental adoption" are both listed candidate
concerns that plainly apply here.

---

## F4 (Low) — The record decides the JSON surface and is silent on the text surface the README advertises

**Anchors**: `0011:C3`, `0011:MVV`, `0011:D-undecided-reporting-shape` (widened
from MVV: every MVV step passes `--as=json`, which sent me to look for the text
path).

Every `0011:MVV` step and every `0011:§testing-strategy` scenario runs
`--as=json`. `0011:C3` constrains Short and Long *help* text but says nothing
about `--as=text` rendering of the new `unknown` `{key, reason}` list.
`0011:D-undecided-reporting-shape` cites terraform's precedent of "reserving the
human-only rendering `(known after apply)` for text mode" — showing the author
was aware there is a text-mode question — and then does not decide it.

`README.md:49-51` states every command accepts `--as text|json` and that text is
"human-readable," so a human caller is a supported audience. The whole product
argument of this RDR (`0011:§approach`, `§risks-and-mitigations`, `0011:F2`) rests
on the reason (`absent` vs `uncomparable`) reaching the caller and telling them
what to do. For the text-mode caller, whether it does is undecided.

**Blocks**: nothing at lock if text mode is deemed out of scope — but the record
should say so. Prevents a text-mode oracle for the reason vocabulary that carries
the entire remedy story.

---

## Not findings (checked, clean)

- `0011:ALT2`'s rejection ("the defect is the default") is a sound product
  argument and is consistent with `0011:§problem-statement` and the Investigation
  prior-art read. No flip-the-default gap.
- `0011:C2`'s `--all` naming survey (`0011:D-naming`) is unusually well grounded
  for a flag decision; no PM objection.
- C3's "a green 0005 suite MUST NOT be cited as evidence" anti-oracle directly
  answers the most likely false-positive path to declaring the outcome delivered.
