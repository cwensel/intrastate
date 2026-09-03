Model: claude-sonnet-5

Delta scope: (1) §key-discoveries "adds no new false-refusal class" claim reversed, cites A7; (2) §consequences new Negative bullet naming the false-refusal class and its re-pricing of the ALT2 economy; (3) D-selection-predicate appended NEW class under A7, open; (4) A2 Status re-scoped to "Verified (scope-limited)", scope note on structural corpus gap; (5) A5 Status Verified→Pending, Method Design Decision→Spike, Falsified note re: cmdbind.go:207; (6) C1 `promise:` clause rewritten as conditional on A6 (Pending); (7) A7 new assumption (Pending, Spike), rarity claim; (8) §approach successor restated as ownership claim, not closure; (9) §mini-check-trace row 3c drops A5 as witness, closing line carries two OPEN rows.

## (a) Internal contradiction check

Searched every section named in the brief for residual claims that the widen refuses nothing new, that the stdin successor closes the axis, or that a description surface exists today.

- §trade-offs / Consequences: carries the SAME text as §key-discoveries' Consequences bullet set (they are the same subsection, reached two ways) — the Negative bullet on position-freedom's false-refusal class is present, correctly hedges "9 of 12 computed shell-free vectors refuse... open (A7)," and correctly notes the stdin axis "waits on the stdin successor" without asserting closure. No stale "refuses nothing new" language survives here.
- §failure-modes: "Silent: a one-word shell string or a stdin-fed interpreter ... lints green — by contract, named in C1." This describes the admission as a documented/named fact, not a closure claim — it does not say the successor withholds or closes the axis. Consistent with the edited §approach.
- §decision-rationale: factor (2) still says freeing the interpreter's position "closes every wrapper with no wrapper table" — but that claim is about the WRAPPER class (argv0 prefix forms: env/nice/timeout/etc.), a different axis than the position-freedom false-refusal class or the stdin axis. Re-read against A7's finding: A7's false-refusal class is a cost of closing the wrapper class, not a wrapper form itself, so this sentence is not contradicted by A7. The Premortem paragraph in the same section still walks through the accepted flag-anywhere class and the promise-wording failure — neither depends on A5 or the reversed key-discoveries claim, and neither restates "refuses nothing new." Clean.
- §testing-strategy: scenario 7 and the preamble ("7 depends on A6 resolving where the description lands") already treat A6 as unresolved — consistent with A6 Pending and with C1's new conditional `promise:` wording. No claim that a description surface exists today.
- 0027:MVV (Minimum Viable Validation) step 5 / §phase-3-promise-wording: "No description surface exists today: `Category` is a bare string and `Categories()` returns identifiers only ... Where it lands is the implementer's call (A6)." This is accurate and matches A6's Pending status — it does not assert a surface exists.
- 0027:S7: "read from wherever Phase 3 lands it" — does not assert a surface exists today; frames it as a future scenario. Consistent.
- Other mini-check tables (`authority`, `disposition`, `oracle`): `authority` and `disposition` don't reference A5, A6, or A7 at all — they describe call-site ownership and lint-outcome routing, which are unaffected by the assumption status changes. `oracle`'s S7 row already says "Today's code: no description surface exists ... so the control starts red by construction" — consistent with A6 Pending, not a settled-fact claim.

No passage found asserting the widen refuses nothing new, that the stdin successor closes the axis as settled fact, or that a description surface exists today. The edits are load-bearing consistently across every section checked.

CLEAR

## (b) Is A7 well-formed?

A7's claim: the false-refusal class "is rare enough in real command bindings to accept as a stated cost rather than bound." That is a RARITY (base-rate) claim.

A7's Evidence section is honest about the gap: "the class is confirmed to EXIST and is trivially constructible... What is NOT established is the base RATE in bindings people actually write — the 12 vectors were selected to exhibit the class, not sampled, so they measure constructibility, not frequency." The Method is Spike, and the spike as executed (12 hand-built vectors, 9/12 refuse) supports only the constructibility half. A7's own Status is Pending, and the Evidence paragraph explicitly proposes the NEXT spike needed to close the rarity claim ("sample real-world `command` bindings... and count how many carry a listed basename as a non-argv0 data word"). This is a well-formed Pending assumption: it does not overclaim Verified on the strength of a constructibility demo, and it names exactly what would verify it.

"If wrong": "the class is common, position-freedom trades a closed bypass for routine false refusals on shell-free tooling, and the 'no wrapper table' economy in §decision-rationale factor (2) is re-priced — ALT2's per-wrapper walk, prohibited by C1's predicate line, becomes the cheaper design." This directly re-prices ALT2 and names the mechanism (§decision-rationale factor 2) — honest, specific, and matches the brief's requirement.

One residual concern, not a defect against the delta-scope gate but worth flagging for Stage 6: A7 is Pending with no verification plan owner/kata named (unlike A5, which got a "Durability caveat" noting no kata tracks it). A7's evidence section names the next spike but doesn't say whether/when it runs before lock. This is the same class of gap as A5's durability caveat, but the brief's "do not re-raise already dispositioned" list doesn't cover it and it isn't one of the four checks — noting only as context, not raising as a new defect since Pending assumptions carrying forward to Stage 6 is the record's stated pattern (A6 does the same).

Method (Spike) is adequate to the claim as currently scoped (Pending, with the frequency spike still to run) — the record does not claim A7 is Verified on constructibility alone.

CLEAR

## (c) Does C1's rewritten `promise:` clause clear the gate?

Gate text (0027:G-assumptions): "no assumption marked `Pending` or `Unverified` may have settled-fact prose elsewhere in the RDR depending on it."

C1's `promise:` clause now reads: "what a reviewer may rely on is the predicate line and nothing more. CONDITIONAL on A6 (Pending): no reviewer-reachable description surface exists in the shipped code today ... IF Phase 3 ships a surface, THESE words are the text it ships. Until A6 settles, this clause promises the predicate line alone, and no admitted-form disclosure may be relied on."

This clears the gate. The clause's operative claim — "what a reviewer may rely on TODAY" — is scoped to the predicate line alone, which is NOT conditional on A6 (the predicate ships regardless). The part that WAS dependent on A6 (whether a description surface carries C1's words) is explicitly bracketed as conditional ("IF Phase 3 ships a surface...") and the clause states outright that until A6 settles, no admitted-form disclosure may be relied on. That is the correct shape: it doesn't assert the future state as settled fact, it names the condition and states the fallback (predicate line only) that holds unconditionally right now. This is a textbook conditional-prose fix, not settled-fact-wearing-a-conditional — the "may be relied on" language is present-tense and true regardless of A6's outcome.

§phase-3-promise-wording: "Where it lands is the implementer's call (A6); that it ships, and carries C1's words, is S7's assertion." This treats A6 as an open implementer decision, not settled fact, and correctly defers the "ships and carries the words" claim to S7 (a test assertion, not a design claim) — consistent with A6 Pending.

0027:MVV step 5: "The description Phase 3 ships is read without provoking a refusal and names both out-of-scope forms in C1's words (S7) — the `promise:` clause's only test; steps 1–4 all pass with no description at all." This describes MVV step 5 as a TEST that S7 will run, not an assertion that the description currently exists or is settled — "steps 1-4 pass with no description at all" is explicitly flagging the gap, reinforcing rather than contradicting A6's Pending status.

0027:S7 (Testing Strategy scenario 7): "This is the `promise:` clause as a test: scenarios 1–5 assert what the predicate does, and none of them would fail if the description were missing, stale, or silent about the admitted forms. That gap is the premortem's second failure ... reaching production with every other test green." S7 is explicitly framed as a not-yet-passing gap-closing test, not a claim that the description exists. Consistent with A6 Pending.

None of the three cross-referenced passages (§phase-3-promise-wording, MVV step 5, S7) assert A6's outcome as settled; all three treat it as an open implementer decision gated by a test. CLEAR.

CLEAR

## (d) Does A5's downgrade leave any dangling citation?

Searched every section pulled for "A5" references.

- §consequences / §trade-offs Negative bullet: "...the second class waits on the stdin successor, which no kata yet tracks (A5)." — cites A5 for the "no kata tracks it" durability fact, which is exactly what A5's new Falsified/Durability-caveat text supports. Not a Verified-witness citation; correctly downgraded framing.
- §approach: "the successor would have to introduce a withholding point, and none exists today (`cmdbind.go:207` sets `cmd.Stdin` unconditionally — A5)." Cites A5 for the falsification fact, consistent with A5's new Pending/Falsified status. Correct — this is citing A5 as the source of the negative finding, not as a Verified closure witness.
- §mini-check-trace row 3c: explicitly edited per the delta list — "A5 is no longer a witness for it (Pending — the successor's withholding premise is falsified), so this row says the form is ADMITTED, not that it is owned." This is the clearest statement in the record that A5 is downgraded, and it is self-consistent.
- §failure-modes / Risks and Mitigations: no A5 citation.
- §decision-rationale: no A5 citation.
- §testing-strategy: no A5 citation (S4's admitted-forms scenario doesn't cite A5).

No passage cites A5 as a Verified witness or relies on its formerly-Verified status. All three surviving citations (§consequences/§trade-offs, §approach, mini-check-trace 3c) point to A5's new content (the falsification finding, the durability gap) rather than its old conclusion.

CLEAR

## New defects

No new defects.

All four delta-scope checks are CLEAR. The nine edits are mutually consistent: the reversed key-discoveries claim, the new Consequences bullet, the D-selection-predicate note, A2's rescoping, A7's new Pending assumption, C1's conditional promise clause, §approach's ownership reframing, and mini-check-trace's two OPEN rows all point at the same underlying finding (position-freedom's new false-refusal class, open under A7; the stdin successor's withholding premise falsified, open under A5's Pending status) without any section left asserting the pre-critique closure claims. A5's downgrade left no dangling Verified-witness citations.
