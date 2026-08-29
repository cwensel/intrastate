Model: claude-opus-5[1m]

# Persona 3 — QA / Tester

Owned set: `S1`–`S5`, `MVV`, `C1`–`C4`. **Widened** to `§implementation-plan`
(Phase 1/3/4), `§failure-modes`, `§risks-and-mitigations`, `A2`, `A3`, `A4`,
and `§illustrative-code` — every widening was forced by a silence in the
owned set that has no line range: an oracle the contracts do not fix
(Q-1, Q-2, Q-4), a phase with no scenario at all (Q-3, Q-5), and an
example that would refuse under its own contract (Q-7).

Verified-clean and NOT reported (checked against source, oracles hold):
S5's wire numbers (13→14 keys, `NumField` 14→15, index 9) match
`decision_table_0010_test.go:416-437` exactly; A2's 28/7 counts reproduce
(`rg -o '"emit":'` = 28, adjacency grep = 7); C4/S5's text rendering
`dispositions: (none)` matches `respond/text.go::flatten`'s empty-map arm
and `decision_table_0010_test.go:699`; `scalar` correctly falls through
both `conformKind` and `conformDomain` (no case), so "never refused" is
mechanically true; escape rules are `sourceRule` too, so S2's "both rule
classes" is decidable pre-`normalizeRules`; `flow next` has no `Emit`
field, so S5's negative leg is trivially green.

---

## Q-1 — C3's "lossless" carry has no ordering oracle for a partitioned domain

**Anchor**: `0024:C3` (claim), `0024:S4` (the test), `0024:C1` (the union clause)
**Severity**: high

C1 fixes the disposition form as "The key's domain is the **union**" of the
per-disposition member arrays, and fixes no order over that union. C3 then
claims the carry is "lossless … the same clause `0002:C22` states for tag
declarations", and the `fidelity` mini-check asserts "value-equality on key,
kind, **domain members**, and each member's disposition — none; the carry is
total". These cannot both hold: taking the union discards both the partition
grouping and the within-partition authored order, and the tag precedent
(`internal/table/source.go:48` → `model.go:134`, `Domain []string`) carries a
single authored array whose order is inherited for free. There is no
authored array for the partitioned form to inherit.

**Test prevented**: S4's "domain members … read identically off the
normalized model". A tester cannot write the assertion — `want` is
undefined. Byte-sorted? Partition-key order then member order?
First-appearance? Each gives a different green. `TestReq146_
EveryEmittedSequenceIsASortedSlice` (`roundtrip_test.go:1332`) makes this
sharper: the package already pins that every emitted sequence is a sorted
slice and that repeated loads produce one value despite randomized map
iteration. A carrier built from a `map[string][]string` domain sub-table
will be **non-deterministic across loads** unless C3 or C1 fixes the union
order, and REQ-146's determinism sweep is the test that will find it —
after implementation, not before.

**Fix shape**: C1 or C3 must state the union's order (sorted is the house
default and satisfies REQ-146), or C3 must drop "lossless"/"total" and say
the carry is member→disposition mapping plus a sorted member list.

---

## Q-2 — `dispositions` has no named carrier surface, so S4 has no read site

**Anchor**: `0024:C3`, `0024:S4`
**Severity**: medium

C3 names the carrier only negatively: "a new model-level type, deliberately
NOT `TagDecl`". No `Model` field name, no type name, no accessor. The tag
analogue is concrete (`Model.Tags map[string]TagDecl`, `model.go:435`), which
is what makes tag-declaration carry testable at all.

**Test prevented**: S4 in full — "key, kind, domain members, and member
dispositions read identically off the normalized model" names no expression
to read them off. Two implementers write two different tests and neither is
wrong. Also blocks C4's join test at the CLI: `flow_resolve.go` holds the
`*table.Model`, but the RDR does not say what it reads from it.

**Also**: A4 licenses reuse of `ConformValue(decl TagDecl, member string)`
(`load.go:793`) — a `TagDecl` parameter — while C3 forbids `TagDecl` as the
carrier. The adapter (emit decl → synthetic `TagDecl` with flattened
`Domain`) is unnamed, and it is exactly where Q-1's ordering choice gets
made a second time. Untested seam.

---

## Q-3 — no scenario covers the JDR 0002 §D1 ECHO/PLAN registration

**Anchor**: `0024:S1`–`0024:S5` (the gap), `0024:C4`; forced widening to
`§implementation-plan` Phase 3
**Severity**: medium

Phase 3 states the obligation explicitly — §D1 requires every new
success-payload field of a projecting verb to be assigned ECHO or PLAN "when
it is added", enforced by a reflective oracle for which "an unassigned new
field is a test failure, never a silent default"; `dispositions` is assigned
PLAN in `docs/jdr/0002-success-envelope-projection.md:73,81`. Phase 3 itself
says the obligation "is not optional and is named here so it is not
discovered as a red reflective test."

**Test prevented**: none of S1–S5 asserts the registration, and the Testing
Strategy's Done criterion is "every scenario below is green" — which is
satisfiable with `dispositions` unregistered. The RDR names a red-test risk
and then omits it from its own acceptance spine. In the 0024-first landing
order there is nothing to run; in the 0023-first order there is a
registration edit no scenario pins.

---

## Q-4 — S3's oracle asserts byte-for-byte equivalence and an appended field at once

**Anchor**: `0024:S3` (with `0024:C2` opt-in leg, `0024:C4`)
**Severity**: medium

S3's Expected reads: "full suite green with the loader change and no model
changed; no new refusal reachable; **the only observable delta is C4's
appended `dispositions: {}`**" — and then "re-runs after the change to prove
**byte-for-byte equivalence**". C4 makes `dispositions` unconditional and
never-omitted, so no resolve payload is byte-identical to its pre-change
form. The two clauses in one scenario contradict.

Compounding it: "full suite green" is unreachable without the licensed
regeneration Phase 3 authorizes — 28 golden literals in
`flow_demand_0011_test.go` go red by construction, plus `TestReq39`'s
`slices.Equal` list and `NumField` check. S3 says "no model changed", which
is true of `models/*.toml` but says nothing about the test corpus, so a
tester reading S3 alone will treat the 28 reds as a failure of S3 rather
than as its licensed cost.

**Test prevented**: the opt-out regression test as written. The tester
cannot decide whether the pass criterion is (a) suite green after the
licensed golden diff and payload delta confined to `dispositions:{}`, or
(b) byte-identity, which is false by C4. A3's Evidence already carries the
correct framing ("the post-change leg (full suite green WITH the loader
change and no model edited)") — S3's "byte-for-byte" is the stray word.

---

## Q-5 — Phase 4 has no scenario and no pass/fail criterion

**Anchor**: `0024:S1`–`0024:S5`, `0024:MVV`; forced widening to
`§implementation-plan` Phase 4
**Severity**: medium

Phase 4 ships two artifacts — a declared `[emit]` block on
`models/examples/pricing-decision-table.toml` and a `dispositions` section in
`docs/cli-output-contract.md`. Neither appears in S1–S5 or in the four MVV
steps (the MVV authors its own throwaway two-rule table, not the pricing
example). Peer RDRs pin exactly this kind of doc obligation with a repo-file
assertion — `escape_shape_0009_test.go:416` and
`flow_rehome_0011_test.go:612` both `readRepoFile(t,
"docs/cli-output-contract.md")` and assert the documented field is present.

**Test prevented**: any assertion that the shipped example lints clean once
declared, and any assertion that the wire contract documents `dispositions`.
Both are silently droppable at implementation with every scenario green.

---

## Q-6 — the three new categories owe a `testdata/neg/` witness that nothing names

**Anchor**: `0024:C1`, `0024:C2` (the categories), `0024:S1`, `0024:S2` (the tests)
**Severity**: medium

S1 and S2 describe inline fixtures, one defect each. The package's standing
convention is one **checked-in** witness fixture per category, asserted by
category: `internal/table/dump_test.go:820-875` (REQ-119) maps every load
category to a `neg/*.toml` file and separately asserts each is present in
`table.Categories()`. The `witnesses` map is hand-maintained and iterated
over itself, so omitting `malformed_emit_declaration`, `unknown_emit_key`
and `emit_value_out_of_domain` **fails nothing** — the convention degrades
silently.

**Test prevented**: the per-category witness test for the three new
categories, and the `Categories()` export check. Neither C2's "join
`0002:C24`'s data-level set" nor any scenario says whether the three slugs
must be added to `table.Categories()` (`category.go:47-71`) — and if they
are not, `loadFindings` still emits them as finding codes while
`Categories()` under-reports the set a consumer can branch on.

---

## Q-7 — the Illustrative Code's `dpa` domain refuses against the model Phase 4 declares

**Anchor**: `0024:§illustrative-code`; bears on `0024:C1`, Phase 4
**Severity**: low

The illustrative block declares `[emit.dpa] domain = ["required", "waived"]`.
`models/examples/pricing-decision-table.toml` authors `dpa = "required"` and
`dpa = "none"` (lines 61, 73, 85, 97). Phase 4 declares that exact model's
emit keys. Copying the illustrative domain refuses `emit_value_out_of_domain`
on two of its four rows.

**Test prevented**: nothing directly — the section is labelled "shapes, not
fixtures". Reported because it is the only worked `[emit]` block in the
record and the nearest thing Phase 4 has to a fixture, and because a domain
that is *plausible but wrong for the real data* is the exact P-1 failure the
Risks section names (a declaration that faithfully copies a mistake).

---

## Q-8 — MVV step 2's "category slug in `code`" is right; the envelope code it rides is unstated

**Anchor**: `0024:MVV` step 2, `0024:C2`, `0024:F1`
**Severity**: low

C2 says the three categories "map to `flow-model-invalid` under `0005:C1`,
and to a blocking finding under `intrastate lint`". The inner finding `code`
does carry the slug (`flow_input.go:185 loadFindings`), so the MVV assertion
is writable. But the **outer** envelope code differs by surface —
`internal/cli/lint.go:183` emits `"model-invalid"` while
`internal/cli/flow_input.go:45` emits `"flow-model-invalid"` — and the MVV
runs `intrastate lint`, i.e. the `model-invalid` branch. The RDR names only
`flow-model-invalid`. A tester asserting the envelope code from C2 writes
the wrong one.

(The `lint.go:97` help text also states `codeModelInvalid` for the
load-refusal branch it does not emit. Pre-existing repo drift, not this
RDR's to fix — flagged only because MVV step 2 is the assertion that trips
on it.)

---

## Q-9 — no scenario pins that `dump` still renders a *declared* key's value raw

**Anchor**: `0024:S4`, `0024:C1`; widened to `§implementation-plan` Phase 1
**Severity**: low

S4's dump leg is "kernel and dump surfaces carry none of it" — i.e. the
declarations are not carried. It does not assert the thing Phase 1 warns
about: `dump.go::renderEmit` (`dump.go:120-132`) deliberately bypasses
`renderValue` because "an emit key … is undeclared, so there is no kind to
consult". Phase 1 rewords that comment. Once a kind exists, rerouting
`renderEmit` through `renderValue` becomes an easy and wrong change —
`renderValue`'s quoting and bracketing key on kind and member count, which
would mutate dump output for declared enum values.

**Test prevented**: a dump-output assertion that a declared enum emit value
renders as the raw authored string, byte-identical to its undeclared
rendering. Nothing in S1–S5 covers it, and C1's "no value is parsed …
canonicalized, or converted anywhere downstream" is the claim that would go
untested.
