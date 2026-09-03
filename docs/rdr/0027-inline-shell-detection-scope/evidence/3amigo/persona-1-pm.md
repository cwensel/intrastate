Model: claude-sonnet-5

# Persona 1 — Product Manager review of RDR 0027

## Findings (severity-ranked)

### 1. (Moderate-High) Promise-wording text — the actual user-facing deliverable — is never verified to ship
**Anchor**: 0027:C1 (`promise:` line), 0027:§phase-3-promise-wording, 0027:§testing-strategy, 0027:MVV

C1's `promise:` clause states the entire point of this RDR in one sentence: "what a reviewer may rely on is the predicate line and nothing more; the lint's user-facing description states the out-of-scope forms in the words above. No such description exists in the shipped code today ... Phase 3 creates it and THESE words are the text it ships." Phase 3 (0027:§phase-3-promise-wording) is a one-line task to write that description.

But every test scenario (Testing Strategy 1-6, MVV steps 1-4, both mini-checks `oracle` and `trace`) asserts on the **category string** and the **refusal detail** (which words matched), never on the **user-facing description text itself** — the thing a reviewer actually reads to learn what is and isn't covered. The premortem in §decision-rationale explicitly worries about exactly this failure mode ("a reviewer who read 'closes the wrapper class' stops reading argv and an `env -S "sh -c …"` string ships... the second is a promise-wording failure") and asserts it's mitigated because "C1's ... line and the lint's user-facing description carry the two admitted forms in those words" — but nothing in the validation plan checks that the shipped description actually contains those words after Phase 3 lands.

**Decision this blocks**: whether the MVV/Testing Strategy as written is sufficient to consider the user outcome (reviewer trust restored via one holdable sentence) actually delivered and verified before Stage 8 implementation is called done. As written, the RDR could ship a fully correct predicate with a stale, wrong, or missing user-facing description and every listed test would still pass.

### 2. (Moderate) The stdin-fed hazard — the one axis proven to actually execute a shell — is deferred to an untracked, unshipped successor
**Anchor**: 0027:A5, 0027:F2, 0027:§background

§Background records that `env -i sh -c 'echo …'` "was verified to actually execute the shell" — this is a live, demonstrated bypass, not a theoretical one. C1 closes the wrapper/argv-position axis but explicitly leaves the stdin-fed axis (`sh -s`, bare `sh`, `python -`, `node -`) admitted, reasoning (A5) that it is "owned by the charted `stdin = "none" | "envelope"` successor." A5's own "Durability caveat" admits: "no kata tracks the successor — it exists only as that `Charted.md` line ... until it ships, C1 names the forms admitted and the reviewer reads argv." So the RDR's answer to one of the two hazards named in the Problem Statement is "the reviewer must go read argv themselves" — which is verbatim the failure mode the Problem Statement opens with ("They discover it only by reading argv themselves — the thing the lint exists to spare them"). The RDR is candid about this (it's a stated, deliberate scope cut), but there is no forcing function (no kata, no tracked follow-on) ensuring the successor ever lands, so the user outcome is only half-delivered indefinitely, by default rather than decision.

**Decision this blocks**: whether shipping C1 alone is an acceptable interim state, or whether landing C1 should be paired with filing/tracking the stdin successor so the second half of the promised outcome doesn't silently rot.

### 3. (Low) Finalization Gate sections are unfilled template placeholders, despite the rest of the record reading as lock-ready
**Anchor**: 0027:§scope-verification, 0027:§contradiction-check, 0027:§assumption-verification, 0027:§proportionality

All four gate-response sections still contain bracketed template instructions rather than written responses (e.g., §scope-verification: "[Confirm the Minimum Viable Validation is in scope and will be executed during implementation, not deferred. State the specific test or proof.]"). This is squarely a product-outcome checkpoint — it's where someone is meant to confirm the MVV proves the promised behavior — and it's empty, while §decision-rationale documents a completed ground-sweep, a fired joint-check dated 2026-08-31, and Status is still `Draft`. This is likely just sequencing (gates fill at a later stage) rather than a defect, but flagging since finding 1 above is exactly the kind of gap Scope Verification should have caught.

**Decision this blocks**: none yet — informational, contingent on Metadata/Status being pre-gate. If Stage 7 (finalize) is imminent, this needs to be filled before lock, and finding 1 should surface there.

## Note on scope widening
I widened past my owned set (§problem-statement, §approach, §decision-rationale, MVV) into 0027:C1, §testing-strategy, §phase-3-promise-wording, §scope-verification/§proportionality/§contradiction-check/§assumption-verification, and A5/F2, because the product question "does the promise actually reach the reviewer" cannot be answered from prose alone — it required checking whether the normative contract's `promise:` clause has a corresponding verification step, which sent me into the contract text and the full validation plan.
