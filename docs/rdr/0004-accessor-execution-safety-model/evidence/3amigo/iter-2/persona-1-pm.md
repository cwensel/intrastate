Model: claude-opus-5[1m]

# Persona 1 — Product Manager (iteration 2, A8 delta scope)

Scope: the A8 re-entry material only — Critical Assumption A8, the read-completeness
normative clause, the "Read completeness" Load-Bearing Decision, MVV / Validation
Scenario 2, the read-truncation Risk and Failure Modes paragraph, the Proportionality
claim, and the mini-check rows carrying read completeness (Disposition Table, Oracle
Discriminability scenario 2, Fidelity Table `read` row, Desk Trace step 2).

Question asked of each passage: does this deliver the user outcome JDR 0001 §D3 opened
the re-entry to secure — "a truncated read makes `exists = false` decide TRUE, and a plan
is produced against state that exists"?

---

## HIGH

### H1. "Read completeness" Load-Bearing Decision — "absent value" resolution defeats the outcome the re-entry exists to secure

Passage: Load-Bearing Decisions, **Read completeness** bullet — "A key the artifact
genuinely does not carry resolves as an absent value"; and its companion in Critical
Assumption A8 ("`<absent>` as a value when the artifact legitimately lacks it") and the
Disposition Table row *Requested key absent from artifact → success → — (absent value)*.

The re-entry was opened because JDR 0001 §D3 named a specific user harm: a presence/absence
predicate deciding on state that was never read. The RDR closes exactly half of that harm.
It makes *unreadability* a refusal — good — but it routes *genuine absence* into the
success branch **as a value**, and the RDR never says what a consumer does with that value.

The kernel is the consumer, and it decides on map presence, not on value content.
`internal/resolve/resolve.go::assemble` writes every `in.Owned` tag into the view keyed by
`t.Key`; `internal/resolve/resolve.go::TagSet.has` returns `ok && provenance match`;
`internal/resolve/resolve.go::missingOwned` fires `owned_state_unavailable` only when
`has` is false. `resolve.Tag` is `{Key, Value string}` — it carries no absence channel.
So a read accessor that faithfully implements this clause and returns a genuinely-absent
key as a value makes that key **present** in the kernel view, and `owned_state_unavailable`
never fires for it. The RDR's own A8 Evidence line cites `missingOwned` — but cites it to
show the kernel cannot distinguish absent from unread, and then hands the kernel a
representation that makes absent look *present*, which is a third and worse state than the
two the JDR was reconciling.

Restated as user outcome: before the re-entry, a genuinely-missing owned tag produced
`owned_state_unavailable` — the user was told what state was missing and who could supply
it (JDR 0001 Principle 4). Under this clause, if the accessor emits an absent-value tag,
that same user gets a resolution computed over a key the artifact does not carry. The RDR
solves the truncation problem and silently regresses the absence problem.

Decision blocked: **what a read accessor returns to the kernel for a genuinely-absent
requested key** — a present tag with a sentinel value, an omitted key, or a typed absence
marker outside `[]Tag`. Implementation cannot pick without either regressing
`owned_state_unavailable` or contradicting this bullet's plain reading. It also blocks
whether `resolve.Input.Owned` needs a shape change (JDR 0001 §D3 option (c), explicitly
left "available if implementation shows the kernel needs to distinguish read-failed from
absent") — the RDR asserts the kernel does *not* need it without testing the absence half.

### H2. "Read completeness" Load-Bearing Decision — the representation escape hatch removes the only thing that would have caught H1

Passage: Load-Bearing Decisions, **Read completeness** bullet, final sentence — "The
contract constrains which branch is taken, not how an absent value is represented in the
success branch; the Resolve spike's `<absent>` sentinel is fixture shorthand, not a
normative representation, and implementation may choose a typed absence marker instead."
Reinforced by the Fidelity Table row *`read` over a requested key set* — "how an absent
value is *represented* is not pinned" — and Oracle Discriminability scenario 2, whose
discriminator is stated purely as *branch* difference.

Read as a PM: the RDR declares the representation question out of contract, then leans on
the spike to prove the assumption. But the spike is where the representation question is
answered — `absentValue = "<absent>"` is a plain `string` in
`evidence/spikes/main.go::read`, sitting in the same `map[string]string` as real values.
The spike therefore proves the branch rule *and* demonstrates the one representation that
breaks the kernel, while the RDR text says representation does not matter. Nothing in the
delta — no normative clause, no MVV assertion, no Disposition row, no Desk Trace assertion
— constrains the success-branch representation, so nothing would fail if implementation
picked the spike's shape.

"Consumers may therefore treat a returned tag set as total over the requested keys" is the
sentence that makes this load-bearing: totality over the requested key set is exactly the
property that makes an absent key look answered. The RDR grants consumers a permission
whose safety depends entirely on a representation it declines to pin.

Decision blocked: whether the read-completeness contract owes the consumer a *distinguishable*
absent value or only a total key set. Until that is pinned, no implementer can write the
kernel-side branch, and Phase 1's "refusal classes" work has no absence type to define.

## MEDIUM

### M1. MVV / Validation Scenario 2 — the acceptance surface tests the truncation half and never asserts the absence half's downstream effect

Passage: Minimum Viable Validation, final sentence — "The read test must assert that an
accessor which can resolve only some of the requested keys takes the refusal branch and is
distinguishable from one whose artifact genuinely lacks those keys"; and Validation
Scenario 2's Expected — "an artifact that genuinely lacks a requested key returns it as an
absent value, not a refusal", proven at `output.txt:13`.

Both assertions stop at the accessor boundary: they check which branch is taken and that
the two fixtures differ. Neither asserts anything about what the absent value *does* to a
consumer. `output.txt:13` (`tags={profile=<absent>,status=Draft}`) is cited as the proof of
the genuine-absence disposition, and it is a proof — of the branch, not of the outcome.

The user outcome the JDR named is a downstream property: a predicate must not decide on
unread state. An MVV that never runs an absent value through a presence/absence predicate
cannot fail when that property breaks. Concretely, the MVV as written passes unchanged
under H1's regression.

Decision blocked: whether MVV Scenario 2 is done at lock. As written it is satisfiable by
the spike as-is, so the implementer has no signal that the absence half is unfinished, and
the Prerequisites checkbox "All Critical Assumptions verified (A1-A8)" is already marked.

### M2. Failure Modes / read-truncation Risk — the silent-failure inventory names one shape where the delta created two

Passage: Failure Modes — "There are two silent-failure shapes, each with a mandatory
guard… A read returns fewer keys than requested and the shortfall reads downstream as
genuine absence — the completeness requirement turns that into an `incomplete_read`
refusal." And the matching Risk bullet — "A read accessor returns a truncated tag set that
a consumer cannot distinguish from genuine absence… Mitigation: Completeness is normative."

The paragraph's own framing — silent failure means a wrong answer nobody is told about —
is the right frame, and it correctly closes truncation-read-as-absence. But the delta
introduced the mirror shape and the inventory does not carry it: **genuine absence read as
present state**, per H1. That shape has no guard named anywhere in the RDR, and by the
paragraph's own standard ("each with a mandatory guard") it should have one or be argued
away.

For a PM this is the diagnosis-quality passage: "Diagnosis starts with the accessor
identity, capability, artifact role, timeout, and expected versus observed tag values" —
none of which help a user whose plan was computed over a key the artifact never carried,
because there is no refusal and no observed/expected pair to inspect.

Decision blocked: whether the absence-as-present shape is a silent failure this RDR guards,
a shape it consciously accepts with recorded rationale, or one it routes to RDR 0007 /
0001. All three are defensible; the RDR picks none, so Phase 2's executor has no obligation
either way.

## LOW

### L1. Proportionality — "part of that same contract" is asserted for the branch rule but is doing work for the representation question too

Passage: Proportionality — "Read completeness is part of that same contract — it is the
success predicate of the read capability, not a separate obligation — so carrying it here
does not widen the RDR."

The argument is sound for the branch rule: which branch a read takes is plainly the read
capability's success predicate, and hosting it here is right. But per H1/H2 the unresolved
question is not the branch — it is what the success branch *hands the kernel*, which is a
seam question between this RDR and RDR 0001's `Input.Owned`. The Proportionality paragraph
covers the first and reads as though it covered the second, which is what lets the
representation escape hatch in H2 look costless.

Decision blocked: whether resolving H1 stays inside RDR 0004 or requires reopening RDR 0001's
`Input`/`Tag` shape (JDR 0001 §D3 option (c)). This paragraph is where that scoping call
would be recorded and it currently forecloses the question by assertion rather than
answering it.

### L2. Desk Trace step 2 — "no CONTRADICTION row" is argued against the wrong pair

Passage: Desk Trace, closing paragraph — "The two assertions that could have collided — 'a
read accessor MUST return typed tag values or a typed refusal' and the completeness
guarantee — are jointly satisfiable because completeness constrains *which* branch the
disjunction takes rather than adding a third branch; step 2's three witnesses exercise both
without conflict."

The pair examined is real and the reasoning about it is correct. But the collision that
matters after the delta is between step 2's absent-value witness
(`tags={profile=<absent>,status=Draft}`, `output.txt:13`) and the kernel behavior A8's own
Evidence line cites (`missingOwned` deciding on map presence) — a cross-layer collision the
Desk Trace's assertion set does not range over, because step 2's assertions stop at "no
mutation / bounded timeout / no direct output".

Decision blocked: nothing on its own — this is a coverage gap in the check that would
otherwise have surfaced H1 before lock. Recorded so the no-contradiction claim is not read
as having cleared the absent-value path.

---

## Not findings (checked, and the RDR gets these right)

- The unreadable-key half of A8 is genuinely delivered: `incomplete_read` is minted at the
  accessor boundary (`evidence/spikes/main.go::read`), witnessed at `output.txt:12`,
  carried as its own Disposition row, its own Fidelity strength, its own Failure Modes
  shape, and its own MVV assertion. This is the JDR §D3 (b) resolution, correctly landed.
- Oracle Discriminability scenario 2's negative control (`output.txt:11`, "an
  implementation that refuses whenever a key is interesting fails this row") is a real
  control against over-refusal, which is the failure mode a completeness rule invites.
- Hosting the rule at the accessor boundary rather than in the kernel is the right call and
  the A8 Evidence line argues it correctly — the kernel provably cannot recover the
  distinction after the fact.
