Model: claude-opus-5[1m]

# Persona 1 — Product Manager: does RDR 0026 deliver the user outcome?

Owned starting set read: `0026:§problem-statement`, `0026:§approach`,
`0026:§decision-rationale`, `0026:MVV`. **Widened** to `0026:C1`,
`0026:§consequences`, `0026:§failure-modes`, `0026:§load-bearing-decisions`,
`0026:§phase-4-surface`, `0026:§testing-strategy`, `0026:§existing-infrastructure-audit`,
`0026:§background`, `0026:§metadata`, `0026:§key-discoveries`. What sent me: the
Problem Statement names two distinct users (a *model author* who declares the
timeout, and an *agent caller* who branches on the refusal) but only argues the
outcome for the second; to find whether the author's outcome lands I had to follow
the diagnostic text out of the contract into the failure modes and the surface phase,
where the outcome either arrives or does not. A silence has no line range, so the
widening was mandatory, not optional.

**Overall**: the primary outcome — *no invocation hangs unboundedly; a class always
arrives* — is delivered, tightly argued, measured, and validated by `0026:MVV`.
The gaps below are all on the **second** user in the Problem Statement (the model
author who must now FIX their tool), whose outcome the record asserts but does not
carry through to a deliverable.

---

## H1 — `0026:F6` × `0026:§problem-statement` × `0026:C1` (`class:`): the record's own headline scenario is the one where the author gets no diagnostic

`0026:§problem-statement` opens on a model author who declared a `timeout`, and
`0026:§background` reproduces the defect with `timeout = "2s"`. That is, by
construction, the case where the context deadline HAS elapsed. `0026:C1` `class:`
then routes exactly that case to `timeout`, and — per 0025:C4, restated in the same
clause — a `timeout` refusal carries **no `Err` and no `Detail`**. `0026:F6` states
this plainly and labels itself "under-explained".

So the entire held-pipe diagnostic apparatus the record builds — the pipe names,
the exit status, the remediation string "close or redirect the helper's inherited
stdio", and the `DrainGrace` constant whose sole stated justification in
`0026:D-the-drain-grace-is-a-mechanism-constant-not-the-rejected-knob` is
"what the grace preserves is the stderr tail … the text that tells a user WHICH
helper held the pipe" — is **unreachable in the record's own motivating scenario**.
It surfaces only on `execution_failure`, i.e. only when the author declared a
timeout generous enough that the child exited first.

The record never states this as a product decision. `0026:§consequences` lists the
`timeout` blindness nowhere; `0026:MVV` step 3 asserts empty `Detail` as *correct*
without noting the author is now told "timeout" for a condition that is not a slow
tool but a broken one. `0026:F6` frames it as a class flip near a deadline "as it
does today for any slow command" — but today's case is a slow command, and this
case is a structurally broken one that no amount of raising the timeout will fix.
The two are not equivalent from the author's seat, and the record treats them as
if they were.

**Decision this blocks**: whether `timeout`'s no-Detail rule (0025:C4) should be
amended by this RDR — which is the RDR that owns amending 0025:C4 — to admit a
held-pipe sub-reason, or whether the author-facing diagnostic is knowingly
sacrificed in the majority case. That is a product call, not an implementation
detail, and it cannot be deferred to Stage 8 because it changes an Override target
this record already claims. If the answer is "sacrificed", `0026:§consequences`
needs a Negative bullet saying so and `0026:F6` needs to stop calling itself
under-explained and become an accepted trade-off with a reason.

Note the sharpened form: this is not the general "class flips near a deadline"
observation. It is that `DrainGrace` — a new shipped constant, one of the two things
this RDR adds to the codebase — is justified in `D-the-drain-grace-…` by a benefit
that the `class:` rule discards in the record's headline case. That justification
and that rule are in tension, and neither cites the other.

---

## H2 — `0026:§phase-4-surface`: the author's remediation path is one unowned sentence with no validation row

Phase 4 reads, in full, that the restated residue should be stated "where authors
read about `timeout` for command entries". No file is named, no doc surface is
named, no owner, no acceptance check. Every other phase in
`0026:§implementation-plan` names a code site; this one names a place in the
abstract. `0026:§testing-strategy` has nine scenarios and none covers the
documentation; `0026:MVV` does not include it; `0026:§prerequisites` does not
gate on it.

`0026:§metadata` **Profile** says `user-facing yes`. A user-facing change whose only
user-facing artifact is an unlocated sentence in the last, smallest phase is the
classic shape of work that gets dropped when the code lands and the tracker closes.
`0026:§consequences` makes two Negative claims that only hold if this phase ships:
that the refusal is "visible, with the pipe named" and that the leaked-process
residue is "explicitly admitted". Admitted *to whom*, and read *where*, is not
answered anywhere in the record.

**Decision this blocks**: whether an author encountering this refusal for the first
time can self-serve. Until Phase 4 names the artifact, the implementer cannot know
what "done" is, and the reviewer at Stage 8 cannot verify it. Name the doc surface
and add it to `0026:MVV` or to the scope gate — otherwise the honest thing is to
delete Phase 4 and record in `0026:§consequences` that the diagnostic lives only in
the refusal string.

---

## M1 — `0026:§consequences` bullet 3 × `0026:F2`: the behaviour change for a *correct* tool is stated as a cost but never sized as a product risk

`0026:§consequences` records that "a tool that answers correctly but leaves a
`setsid` daemon holding its stdout is now refused rather than accepted". `0026:F2`
carries the engineering side well — A5's 160 trials, 27.7 ms worst case, load
sensitivity, ~18x tail headroom. But that measurement addresses the *in-group*
helper that merely misses the bound. It does **not** address the population this
bullet is about: tools that deliberately and permanently escape the group, answer
correctly, and are today accepted.

Those authors experience this RDR purely as a regression — their working reader
starts refusing. The record has no estimate, no survey, and no migration note for
them; `0026:§risks-and-mitigations` bullet 3 offers "the refusal names the held pipe
and the bound, so the author sees exactly what to fix", which presumes the author
can change the helper. If the helper is a third-party binary that daemonizes, they
cannot.

The RDR is, I think, still right to refuse — `0026:§decision-rationale`'s correctness-fit
row is convincing that parsing output from a stream with an unreachable live writer is
worse. But "right" and "unsurveyed" are different states, and the record currently
presents an unsurveyed population as a one-line accepted negative.

**Decision this blocks**: whether this ships as a plain behaviour change or needs a
migration note / release-note callout. Cheap to close: one sentence in
`0026:§consequences` stating whether any known reader in the repo's own fixtures or
the corpus exhibits the escaping-but-correct shape, and if unknown, saying it is
unknown.

---

## M2 — `0026:§problem-statement` names the agent caller's outcome; nothing in the record validates it end to end

The stated outcome is that "a caller (an agent driving `flow resolve` in a pipeline)
waits forever with nothing to branch on" stops being true. `0026:MVV` and every
`0026:S1`–`0026:S9` scenario assert at the binding/executor level — elapsed time,
class, `Applied()`, parsed-or-not. `0026:S1` does reach for `flow_exec.go::accessorFailureOf`
rendering, which is the closest the matrix gets. But no row asserts the thing the
Problem Statement promises: that `flow resolve` itself terminates and emits a
branchable result.

For most RDRs I would not raise this — the executor assertion is a fair proxy. I
raise it here because the defect being fixed is *specifically* a whole-CLI hang
(`0026:§background`: "still blocked after 15 s"), and the reproduction was done at
the `flow resolve` level, not the binding level. The validation has moved a layer
down from where the bug was observed.

**Decision this blocks**: whether Stage 8 can close `intrastate#936e` on the binding
tests alone. Adding "and `flow resolve` exits non-zero with the class rendered"
to `0026:MVV` step 3 is a small edit that makes the fix testable against the
original report.

---

## L1 — `0026:§approach` and `0026:§decision-rationale` argue engineering precedence; neither states the user-visible promise in one line

`0026:§approach` opens on "when a bounded wait and a whole read conflict, the bound
wins" and `0026:§decision-rationale`'s deciding rows are "correctness fit" and
"prior-art alignment". Both are the right frames for the engineer. Neither section
contains a sentence an author could read as the promise being made to them — the
nearest is `0026:§consequences` bullet 1.

This is a readability note, not a defect: the outcome IS derivable, and `0026:C1`
`precedence:` is unambiguous. But `0026:§approach` is the section a future reader
lands on first, and it currently requires them to reconstruct the user promise from
mechanism.

**Decision this blocks**: nothing hard. It affects whether a later reader
re-litigates the precedence because they cannot see what it bought.

---

## What is NOT a finding

- The `timeout + 2·WaitDelay` bound is measured (`0026:A6`, 418-487 ms headroom),
  the legs' non-independence is disclosed rather than hidden in `0026:C1` `bound:`,
  and `0026:S8` tests the competing-legs case. This is unusually honest and I have
  no objection.
- `0026:§decision-rationale`'s rejection of Alternative B (accept partial output) is
  the load-bearing product call and it is correctly reasoned from the user's seat: a
  write's read-back confirming a write against unproven output is the outcome that
  would actually harm someone.
- The unfilled `0026:G-cross-cutting` / `0026:G-proportionality` / `0026:G-scope`
  gate bodies are template text; `0026:§metadata` Status is `Draft`, pre-lock, so
  these are expected to be answered at Stage 7 and are not raised as findings here.

---

## Severity summary

| # | Anchor | Severity | Blocks |
| --- | --- | --- | --- |
| H1 | `0026:F6`, `0026:C1` `class:`, `0026:D-the-drain-grace-is-a-mechanism-constant-not-the-rejected-knob` | High | whether `timeout`'s no-Detail rule is amended or the diagnostic is knowingly sacrificed in the headline case |
| H2 | `0026:§phase-4-surface` | High | what "done" means for the only user-facing artifact; Stage 8 verifiability |
| M1 | `0026:§consequences`, `0026:F2` | Medium | plain behaviour change vs. migration note for correct-but-escaping tools |
| M2 | `0026:§problem-statement`, `0026:MVV` | Medium | whether `intrastate#936e` can close on binding-level tests alone |
| L1 | `0026:§approach`, `0026:§decision-rationale` | Low | future re-litigation of a settled precedence |
