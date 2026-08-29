Model: claude-opus-5[1m]

# Persona 1 — Product Manager: does this RDR deliver the user outcome?

Owned set read: `0023:§problem-statement`, `0023:§approach`,
`0023:§decision-rationale`, `0023:MVV`.

**Widened.** Three things sent me out of the owned set:

1. `§problem-statement` states the outcome in **transcript tokens** but
   every measurement in the record is in **bytes** — the conversion is
   nowhere in the owned set, so I went to `§performance-expectations`,
   `§consequences`, and `evidence/spikes/a1-byte-width.md` looking for it.
   It is not there either. That silence has no line range, which is why
   it is F1 below.
2. `§problem-statement`'s headline number (1529 B) is not the number
   `0023:A1` verified (708 B), so I read `A1` and its spike to see which
   one the reader is meant to believe.
3. `§approach` promises the benefit to "chained or secondary calls" but
   names no adopting consumer; the only named consumer is in
   `§background` (a navigator skill) and there is no adoption step in
   `§implementation-plan`. That sent me to Phases 1–3 and
   `§prerequisites`.

Findings are severity-ranked. Severity here is *PM* severity — does it
put the stated user outcome at risk, or leave a decision unmakeable —
not implementation risk.

---

## HIGH

### H1 — `0023:§problem-statement` (with `0023:A1`, `0023:§consequences`): the outcome is denominated in tokens, but nothing in the record measures or estimates tokens

The Problem Statement makes the outcome unambiguous and correct: the
caller "pays for every output byte as transcript tokens re-paid on each
subsequent turn," and explicitly demotes wall time ("wall time is 5ms").
Tokens are the currency. But every downstream measurement — `0023:A1`'s
verified table, `§performance-expectations`, `0023:MVV` step 6,
`evidence/spikes/a1-byte-width.md` — is denominated in **bytes**, and no
passage anywhere states the bytes→tokens relationship for this payload
class, nor claims the two are proportional.

This matters for the outcome, not just for tidiness: JSON payloads of the
shape measured here are dominated by punctuation, quotes and short repeated
keys, which tokenize very differently from prose. A 79.4% byte reduction
on a 48-key `observed` map is not self-evidently a 79.4% token reduction,
and the record never says whether the author checked.

**Decision it blocks:** whether the flag clears the bar the RDR itself
sets. `0023:A1`'s "If wrong" clause routes the whole RDR back to "no
projection axis" if the saving is too small to justify a new surface on a
locked verb — but that bar is a token bar and the verification is a byte
verification, so the go/no-go test as written cannot be failed by the
evidence as gathered. It also blocks sizing the ask against the cost
(`§consequences` lists three negatives, including a standing partition
obligation on every future field): a reader cannot weigh "worth a
permanent doctrine" against a saving in the wrong unit.

Cheapest fix consistent with the record's own standard: one sentence in
`§performance-expectations` or `0023:A1` either (a) stating a measured
token count for at least the 48-fact class, or (b) stating explicitly
that bytes are used as a proxy and why that is sound for this payload
shape — so the substitution is a recorded decision rather than a silent
one.

### H2 — `0023:§problem-statement` vs `0023:A1`: the headline motivating number (1529 B) is contradicted by the verified evidence (708 B) and is never retired

`§problem-statement` opens with "Measured on a 48-fact decision-table
call: full JSON is 1529 bytes, the `emit` + `rule` answer ~430 — about
72%". `0023:A1` was raised precisely to test that transfer, and its
Verified evidence measures the comparable 48-fact synthetic at **708 B
default → 146 B projected**. `§consequences` and `§key-discoveries`
both correctly explain that the *ratio* moved (72% echo share vs 79.4%
shipped reduction), but **no passage reconciles the absolute width**: the
Problem Statement's full payload is more than double A1's. The spike
constructed its own synthetic model rather than measuring the original
call, so the 1529 B call is now unattested and unretired.

The prose reads as if only the percentage was corrected. A reader
arriving at `§problem-statement` — the section that establishes the size
of the user pain — still takes 1529 B as the measured baseline, and
nothing downstream tells them it isn't.

**Decision it blocks:** the proportionality judgment, and therefore
`§proportionality` (`0023:G-proportionality`, currently an unanswered
template block). Is the motivating call a ~1.5 KB payload or a ~0.7 KB
payload? A 562-byte absolute saving on a 5ms call is a materially
different product case from a 1100-byte one, and it is the case a
reviewer must accept to justify a permanent new flag plus a standing
cross-RDR partition doctrine. Also blocks a clean read of `0023:A1`'s
"If wrong" route, since the reader cannot tell whether A1 confirmed the
seed or quietly replaced it.

---

## MEDIUM

### M1 — `0023:§approach` / `0023:MVV` / `§implementation-plan` Phases 1–3: nothing in the plan puts the flag in a caller's hands, so the RDR ships a capability and calls it the outcome

The user outcome as stated in `§problem-statement` is that *an
agent-driven caller stops re-paying the echo*. The three implementation
phases deliver: the flag on the verb (Phase 1), the oracle battery
(Phase 2), and docs + extended help (Phase 3). `0023:MVV`'s six steps are
all CLI-level assertions, and its End-state is "one flag, two widths, one
decision" — a capability statement, not an outcome statement.

`§background` names the actual originating consumer ("filed while
shrinking a consumer's navigator skill") and even names the sequencing
need ("land before or with the consumer migrations that add more models
per skill run"). But no phase, no prerequisite, and no MVV step touches a
consumer. The record therefore does not say who adopts `--plan-only`,
when, or how anyone would know the token saving was realised. Since the
flag is opt-in and default-preserving (`0023:C1`), a landing with zero
adopters produces exactly zero of the benefit the Problem Statement
motivates — and every oracle stays green.

I am not asking for consumer work to be pulled into scope; the opt-in
design is right. The gap is that the record never states the boundary —
that delivering the *outcome* requires a caller-side change owned
elsewhere, and names where.

**Decision it blocks:** whether this RDR can be called done. Right now
"done" is defined purely in `§testing-strategy` terms ("done = every
oracle green and the MVV run recorded"), which is satisfiable with the
user outcome entirely undelivered. It also blocks the sequencing
decision `§background` gestures at: "land before or with the consumer
migrations" is guidance, but `§prerequisites` lists only assumption
verification and the 0024 ordering tolerance, so nothing makes the
sequencing binding or even visible at landing time.

### M2 — `0023:§problem-statement` (last two sentences) states the routable question in terms the chosen solution does not fully answer, and `§decision-rationale` does not close the gap

The Problem Statement's routable question has two halves: "should the
resolve payload have a caller-controlled projection axis at all — and if
so, does it live on the flag surface... in the `respond` gateway... or
nowhere?" It then adds a third: "Each answer also decides whether future
verbs inherit the projection."

`§decision-rationale`'s matrix scores four options against the first two
halves and answers them cleanly (O1 wins). The third half — do future
verbs inherit — is answered only obliquely: the "Extensibility" row says
O1 gives "one rule: project the verb result pre-`respond.OK`, echo group
only," and `§decision-rationale`'s joint-check defers the general
question to `JDR 0002 §D1`'s "later-verb clause," noting that clause
"obliges [0021] to nothing." So the Problem Statement promises this RDR
decides inheritance, and the body decides that it doesn't decide it.

That may well be the right answer — but it is a different answer from the
one the Problem Statement told the reader to expect, and neither section
says so.

**Decision it blocks:** a future verb author's first-hour question — "am
I obliged to offer `--plan-only`, permitted to, or forbidden from doing
so without a new RDR?" As written the answer requires reading a JDR the
RDR only cites. One clause in `§approach` or `§decision-rationale`
stating the actual disposition (permitted, obliged to nothing, no
inheritance minted here) would retire the promise the Problem Statement
made.

---

## LOW

### L1 — `0023:MVV` step 6 and `0023:A1`: the MVV re-records the byte table but names no threshold, so the re-measurement cannot fail

MVV step 6 says "Record default vs projected byte counts on the
motivating-model shape and one gate/write-heavy fixture (A1's table...)".
"Record" is the whole obligation — no floor is stated. `0023:A1`'s "If
wrong" branch is the record's only stated failure condition ("saves too
little to justify a new surface") and it too is qualitative ("materially
smaller"). At implementation time, a table showing a 5% saving would
satisfy MVV step 6 as literally written.

This is LOW rather than HIGH because A1 is already Verified with a large
margin (42–79%), so the realistic risk is small. But it is the same
defect as H1 in miniature: the record's economic bar is never made
falsifiable.

**Decision it blocks:** the `§scope-verification` gate response
(`0023:G-scope`, presently an unanswered template block), which asks the
author to "state the specific test or proof" that the MVV is executed and
not deferred. Step 6 is a recording, not a test, and the gate response
will have to say so.

### L2 — `0023:§consequences` "Scope of the guarantee" bullet: the refusal-path carve-out is stated as a design fact, not sized against the caller's actual traffic

`§consequences` correctly and honestly scopes the benefit: "the token
reduction covers successes only... so a probing caller that mostly
refuses saves nothing, by design" (grounded in `0023:A6`). Good
disclosure. But the record nowhere indicates what fraction of the
motivating consumer's calls are refusals. The Problem Statement's
motivating pattern — "the second `--outcome` call in a chain receives the
same 48 facts again" — is a success-path chain, so the carve-out is
probably immaterial; the record just never says so.

**Decision it blocks:** nothing hard. It slightly weakens the expected-
value case a reviewer would build from `§consequences` + `A1`, since the
realized saving is (byte saving) × (success rate) and only the first
factor is known. One clause noting that the motivating chained pattern is
success-dominated would close it.

---

## Explicitly checked and NOT flagged

- `0023:§approach`'s choice of a fixed partition over a caller-supplied
  field list is well motivated from the user's side (a caller can never
  drop `escaped`), and `§decision-rationale`'s matrix supports it. No
  finding.
- `0023:MVV` steps 1–5 are concrete, executable, and each maps to a
  stated contract. No finding.
- The opt-in default (`0023:C1`, `0023:BR1`) is the correct product
  choice and is defended on its own ground rather than by citation
  alone. No finding.
- `§consequences`'s three negatives are stated plainly rather than
  minimized. No finding — this is the section doing its job.
