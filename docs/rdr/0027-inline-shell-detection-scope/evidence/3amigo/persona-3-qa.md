Model: claude-sonnet-5

# Persona 3 — QA / Tester findings on RDR 0027

Overall: this record is unusually disciplined for testability. C1 states an explicit
predicate, a report contract (which two words get named), and an explicit
out-of-scope line. §testing-strategy's six scenarios each carry an oracle, and the
Implementation Plan's `oracle` and `trace` mini-checks go further than most records by
naming their own weak oracle (Scenario 4 / MVV step 3) and pre-empting it with a
discriminating control. Findings below are the residual gaps only.

## Findings (severity-ranked)

### 1. (Medium) Scenario 4's stated oracle is weaker than the oracle the record itself says is required — the strengthening lives in a different section and is not restated where the test would actually be written

**Anchor**: 0027:S4 (§testing-strategy), vs. 0027:MVV (mini-check `oracle`, under
§minimum-viable-validation)

**Test I cannot write from §testing-strategy alone**: "assert the admitted forms
(`sh ./gate.sh`, `env -S "sh -c echo"`, `sh -s`, `sh`, `sh -es`, `python -`, `node -`)
lint green, in a way that actually discriminates a no-op predicate from a correct
one."

**Why**: Scenario 4's "Expected" is only "every one lints GREEN." The record's own
`oracle` mini-check table (fired explicitly "by MVV step 3 and Testing Strategy
scenario 4") calls this exact oracle out as the weak row: "an absence-of-error
assertion that a no-op implementation also satisfies," and prescribes a fix — "Control
= assert green **and** that the binding loads with the argv unchanged, not merely that
no error was returned." That fix is never folded back into Scenario 4's own Expected
line. §testing-strategy states its own done-criterion as "every scenario below green"
— an implementer or test-writer working from that section (the section this record
designates as the test plan, driven through `table.Load`) has no textual instruction
to add the argv-unchanged assertion; they would have to independently notice the
mini-check table lives in a different top-level section (Implementation Plan, not
Validation) and back-port its stronger oracle. As written, a test satisfying Scenario
4's literal text (assert `lint` returns no error for the seven vectors) would pass
against a stub that always returns "green," and nothing in §testing-strategy flags
that as insufficient — only the MVV mini-check does.

**Fix shape** (not prescribing wording): restate the argv-unchanged control inside
Scenario 4's own "Expected" line, or add an explicit cross-reference from Scenario 4 to
the `oracle` mini-check so the strengthened assertion travels with the scenario a test
author will actually open.

### 2. (Low) Scenario 3's oracle omits the detail-text assertion that Scenario 1 and C1's `report:` clause both require

**Anchor**: 0027:S3 (§testing-strategy), vs. 0027:C1 (`report:` line), vs. 0027:S1

**Test I cannot write with a stated pass/fail bar**: "for each of the nine `env`-option
forms (`-i`, `-u FOO`, `--unset=FOO`, `-0`, `-C DIR`, `--chdir=DIR`, `--`, the mixed
chain), assert the refusal detail names the correct matched words."

**Why**: Scenario 3's Expected is only "all refuse `command_shell_interpreter`." C1's
`report:` clause is explicit that "the detail additionally names the two matched
words," and Scenario 1 operationalizes that ("detail naming the matched words `sh
-c`"). Scenario 3 doesn't carry the same detail-text bar, so a test written strictly to
Scenario 3's text would assert category only — which the record itself flags elsewhere
(mini-check `oracle`, S5 row: "a bare category assertion passes even with today's
message") as an insufficient oracle for exactly this contract. This is lower severity
than Finding 1 because Scenario 1 already establishes the pattern in the same section
and a competent test author would likely carry it forward by analogy — but the
scenario as literally written doesn't require it.

### 3. (Low / informational — widened outside the owned set) The interpreter deny-list's actual membership is never shown in 0027; scenarios exercise only `sh`, `python`, `ruby`, `node` as "listed" interpreters

**Anchor**: 0027:C1 ("interpreter set: OPEN (deny-listed, not closed), unchanged from
0025:C5"), 0027:S1, 0027:S4, 0027:S5

**Widened because**: to check whether "every listed interpreter word" (C1's predicate)
is testable, I need the list; 0027 doesn't define it and only cites 0025:C5 as owner,
so I looked at how many distinct interpreter basenames the scenarios actually cover.

**Test I cannot write from 0027 alone**: "assert the predicate refuses every member of
the deny-list when paired with its own inline-code flag" — a full-coverage test over
the list.

**Why**: this is expected and probably fine — 0027 explicitly scopes itself to argv
*position*, not list *membership*, and states the list is owned and grown by amendment
of 0025:C5, not this record. I flag it only because two scenarios lean on specific
list-membership facts asserted in prose rather than in a normative table: S4/trace row
3c states "python IS on the deny-list" as load-bearing for why the out-of-scope line
must be channel-scoped rather than spelling-scoped, and S5 relies on `ruby` being
listed. Neither 0027 nor the excerpt available to this persona shows the list itself,
so a QA reader cannot independently verify those two "is a member" premises without
opening 0025:C5. Not a blocking gap for 0027's own contract (which doesn't own the
list), but worth a forward pointer (e.g., "list is 0025:C5 line N") so a test author
doesn't have to hunt for it when writing S4/S5.

## Not a finding

- F1/F2/F3 (Failure Modes) are each observable and each map onto a scenario: F1
  (visible refusal) ↔ S1/S5's detail-text checks; F2 (silent admission) ↔ S4's
  admitted-forms list; F3 (no recovery state) is a non-claim and needs no test. No gap
  here.
- 0027:MVV and 0027:§testing-strategy are near-duplicates (same six scenarios,
  restated), which is a redundancy/authoring concern, not a testability gap — every
  scenario in both carries the same stated oracle, so QA coverage is identical in both
  places.
- The `trace` mini-check's per-step witness column (spike references A1-A5) gives every
  MVV step an independent falsifiability check beyond "did the assertion pass" — this
  is stronger than most records' validation sections and not flagged as a gap.
