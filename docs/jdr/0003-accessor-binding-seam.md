---
authors: Chris K Wensel <cwensel@retrofit.sh>
state: settled
cluster: 0004, 0025, 0026, 0027, 0028
labels: intrastate, accessor-binding, refusal-carrier, execution-failure, rdr-cluster
---

# JDR 0003 What Does an Accessor Binding Report When It Refuses, and Where Does the Reason Ride?

*Go identifiers and field names below are **non-normative**; the RDRs own
exact contracts. Source: the 0026×0028 Stage-2 joint check, 2026-08-31.*

## Problem statement

Two drafts are about to mint new producers of one shipped refusal class,
`execution_failure` (`internal/accessor/model.go::ClassExecutionFailure`), and
each fences its own carrier for the sub-reason a caller — an agent loop —
needs in order to tell *not applied / did not run* from *ran, output
unproven*. 0026 puts the held-pipe reason (pipes held, bound, the child's exit
status, remediation) in `ExecError.Err` prose and the stderr tail in `Detail`;
0028 puts a rule id, a reason token and `applied: false` in `Detail` and
asserts no typed reason field is needed. One discriminator, two conventions,
no owner: 0004 owns the class (Implemented; silent on sub-reasons), 0025 owns
the `*accessor.ExecError{Detail, Err}` carrier (Implemented), and 0026/0028
are the producers. Answered per-RDR the carrier drifts and the 7.1 gate finds
the contradiction after lock. This registry is the home for what a binding
reports across the accessor-binding seam when it refuses. 0027 amends the
same declaration seam (`cli/0025:C5`) and is bound as a member; no entry
lands in it today.

## Principles

1. **One home per contract; cite, never restate** — the engine doctrine
   (`$RDR_HOME/stages/README.md`, *Single-source each contract*); why this
   registry exists.
2. **Refuse, never guess** — `0004:C14`: a post-run refusal "MUST NOT be
   represented to callers as a write that did not occur"; the LBD: "guessing
   absence is the failure this contract exists to prevent".
3. **`Detail` faces the caller; `Err` faces the executor** — `cli/0025:C4`:
   "`Detail` is the 4 KiB stderr tail, not `Err.Error()`", and `Err` wraps
   the offending error "so `errors.Is` … survives the trip to the refusal
   site" — the refusal site, not the envelope.
4. **The class set is closed and exit-mapped** — `RefusalClasses()` is "the
   closed accessor refusal-class set"; JDR 0001 §D10: "one CLI code per
   caller-branchable failure; inner discriminators ride one structured field".
5. **The applied sense is the executor's** — shipped `Refusal.applied` is
   unexported "because only the write path — which knows whether the command
   already ran — may set it" (`Refusal.Applied()`, `0004:C14`).

## D1 — Where does an `execution_failure` sub-reason ride?

**Constraint** — shipped code and P2–P5 determine it; nothing to negotiate.

- **(a) Prose per producer** — `Err` text (0026) or a `Detail` token
  convention (0028): the two drafts as written. Ruled out by shipped code:
  `executor.go::refusalOf` copies only `ExecError.Detail` onto the refusal
  and nothing downstream inspects `Err`, so 0026's "so a caller can tell"
  reason never reaches a caller; and `applied: false` in `Detail` is a second
  copy of `Applied()`, which is already false for every `execution_failure`
  by construction (P1, P3, P5).
- **(b) A typed discriminator every producer sets** — a tri-state or a reason
  enum on `ExecError`/`Refusal`. A reason enum is a second closed set beside
  the class set with no caller branch to justify it (0028's four anchor
  reasons share one remedy: fix the model); a producer-set applied flag
  inverts P5. What survives of (b) is executor-side: the executor needs one
  typed fact from the binding to set `applied` where its "an `Apply` error
  precedes mutation" assumption (`executor.go::Write`) is false.
- **(c) Split the class** — a not-applied class for writes. Ruled out by P4,
  `0004`'s closed set, and `0028:C1.3`'s own "no new refusal class"; it also
  splits on the wrong axis — the caller's branch is the applied sense, which
  already rides `Applied()`.

**Resolved:** three rules, no new field, no new class.

1. **Caller-facing sub-reason text rides `Detail`, the one slot.** `Detail`
   is the binding's bounded caller-facing diagnostic; a command binding's
   instance is the stderr tail (`cli/0025:C4`'s spelling stands for command
   bindings, unedited). Composition order is the shipped one: applied-sense
   text first, producer text, then the stderr tail.
2. **`Err` is typed, executor-facing, never caller text.** A producer that
   needs the executor to act wraps a typed error in `Err` for `errors.As` at
   the refusal site.
3. **The applied sense is `Applied()`, set by the executor, never prose.**
   A held pipe on the write path after the direct child was reaped is a
   post-run refusal in `0004:C14`'s sense: class stays `execution_failure`,
   `Applied()` is true, and the CLI's "may have been applied" text keys on
   `Applied()` (the shipped mapping keys on class and phase, which this case
   slips past).

Lands in **0026** — `cli/0026:C1`'s `refusal:` line: the held-pipe reason
(pipes, bound, exit status, remediation) moves from `Err` prose into `Detail`
ahead of the stderr tail; `Err` wraps a typed held-pipe error (its instance
names it); the executor's write arm and the CLI's applied-sense rendering
gain the `Applied()` derivation (rule 3). Lands in **0028** — `0028:C1.3`'s
`order:` line drops `applied: false` from `Detail` (rule id and reason token
remain); `0028:A8` is confirmed as a constraint of this entry, and its
If-wrong "owned here" is void — a typed field on the carrier is this
registry's, not 0028's. **0004** and **0025** are Implemented and not edited;
0026 and 0028 are the successors at their seams and cite this entry.
**0027** is bound as a member; nothing lands in it.

## D2 — Is an argv word bound at invocation inside the inline-shell promise?

**Decided** — (a), 2026-09-03; spans 0027, 0028. Filed and decided at the
7.1 cluster reconcile (evidence and research:
`docs/rdr/cluster-reconcile/0026-0027-0028/`, `research-d2-d4.md`).

`cli/0027:C1` promises a reviewer the predicate line plus two admitted forms
named by name — a one-word shell string and a stdin-fed interpreter — over
"argv WORDS only". `cli/0028:C1.6` makes `{tag.<key>}` a whole argv element
bound from `--tag` at invocation, at any position, with no value rule:
`["{tag.a}","{tag.b}","{tag.c}"]` lints green under both records and executes
whatever the caller binds. Shipped code decides nothing here — `carrierDefect`
clause 3 tests placeholder membership only, and `cmdbind.go::substitute`
refuses a `-`-prefixed `{artifact}` value while nothing refuses a
`-`-prefixed tag value. Neither record owns the answer: 0027 owns the
promise's wording, 0028 the placeholder's admission.

- **(a) A rule on the placeholder (0028)** — `{tag.<key>}` never at argv0,
  and a bound value beginning with `-` refuses before spawn, mirroring
  `substitute`'s `{artifact}` rule; the interpreter deny-list stays static.
  Amends `0028:C1.6`'s "no other change" line — a 0028 re-lock, contract-
  scoped.
- **(b) Disclosure in the promise (0027)** — `0027:C1`'s out-of-scope list
  names "an argv word bound at invocation" as a third admitted form and the
  `--help-all` text says so. Amends `0027:C1` — a 0027 re-lock; the hole
  stays, disclosed.
- **(c) Both.**

**Resolved: (a).** The two concrete evasions (an interpreter at argv0, a
flag-shaped value) close under a static check, and the promise stays true
without naming a form that exists only to be admitted. Stake: a reviewer
approves a model on the strength of the shipped promise while the executed
argv is a string they never saw. 0028's Cross-Cutting gate response ("adds no
new source of caller-supplied value") is false under C1.6 whichever way this
resolves; the answer supersedes it. Grounded: Go `os/exec` runs no shell,
so the class is argument injection (CWE-88), not shell injection; Go `flag`
and git treat a leading dash as flag-shaped and offer `--` only where every
child honours it; `cmdbind.go::substitute` already refuses a `-`-prefixed
`{artifact}` rather than inserting `--`. Lands in **0028** at `C1.6`
(admission and binding) — a contract-scoped re-entry; **0027** unchanged.

## D3 — Which exit group do 0028's stale-model refusals take?

**Decided** — (b), 2026-09-03; spans 0028 and JDR 0001 §D10's code table
(landing RDR 0005, Implemented). Filed and decided at the 7.1 cluster
reconcile.

`0028:C1.3` mints `edit_anchor_unmatched` / `_ambiguous` / `_collision` /
`_unstable`, `edit_clear_undeclared` and an unreadable target as
`execution_failure` with "no new refusal class". JDR 0001 §D10 rule 1 maps
execution failure to exit 3 — "repair the environment and re-run the same
request unchanged" — and `docs/cli-output-contract.md` says a refusal about
the request must never exit 3 or the caller spins. §D10 rule 7 already pulls
read-back mismatch out to exit 2 on exactly that reasoning. D1 deferred this
question "to 0028 by citation"; 0028 cites §D10 nowhere.

- **(a) Exit 3 as written** — class mapping unchanged; the remedy rides
  `Detail` prose. Cost: a retry-on-exit-3 agent loops on an artifact it must
  instead fix.
- **(b) Exit 2, no new class** — the executor-facing typed `Err` (D1 rule 2)
  discriminates at `flow_exec.go::accessorFailureOf`; a distinct CLI code
  (spelling non-normative, e.g. `flow-write-anchor-stale`) carries the rule id
  and reason token in `findings[]` per §D10 rules 2–3. Cost: one code added
  to 0005's table by citation repair; 0028's class set and `Detail` unchanged;
  touches the CLI arm 0026 also re-keys (D4).
- **(c) A new accessor refusal class** — ruled out by 0004's closed set and
  `0028:C1.3`.

**Resolved: (b).** Stake: an agent that branches on exit codes retries a
stale anchor until its budget expires. Grounded: sysexits separates
`EX_TEMPFAIL` (retry later) from `EX_DATAERR`; HTTP 412 is the client-must-act
precedent for "the target changed underneath you"; this org's `cli/0028`
`translateDDLError` and `cli/0120` A1 both route a stale anchor to the exit-2
group as a loud refusal. Lands in **0028** (the typed `Err`, the CLI arm, a
citation of this entry at `C1.3`) and JDR 0001 §D10's table (one code, by
citation repair).

## D4 — How does the applied sense reach the CLI envelope?

**Decided** — (b), 2026-09-03; spans 0026, 0028, and JDR 0001 §D10. Filed
and decided at the 7.1 cluster reconcile.

D1 rule 3 makes `Applied()` the discriminator at the accessor seam, and
JD-1's stake is an agent loop deciding whether to re-read before retrying a
write. At the envelope nothing carries it: no `json:"applied"` field exists,
`Applied()` has no non-test caller, and 0026's applied held-pipe write refusal
and 0028's not-applied pre-mutation refusals both render as
`flow-accessor-failed` / exit 3, distinguished only by whether `detail`
carries the "may have been applied" sentence (`0026:S1` asserts that prose;
`0028:S27` asserts nothing at the envelope). §D10 rule 2: one CLI code per
caller-branchable failure. §D10 rule 7's precedent: read-back-incomplete and
read-back-timeout are distinct exit-3 codes whose message carries the applied
text.

- **(a) Prose** — `detailMayHaveApplied` in `detail`, code unchanged; `0026:S1`
  as written. Cost: callers string-match `detail` to branch.
- **(b) A distinct exit-3 CLI code for the applied write refusal** (spelling
  non-normative, e.g. `flow-write-failed-applied`), keyed on `Applied()` in
  `accessorFailureOf`'s write arm — the concrete form of `0026:A8`'s "gain its
  `Applied()` key"; `detail` text unchanged. Cost: one code in 0005's table.
  Contradicts no fence in `cli/0026:C1`; 0028's refusals keep
  `flow-accessor-failed`.
- **(c) A structured field** — `applied` in `findings[]`, §D10 rule 3's slot.
  Cost: a findings entry on an accessor refusal, a new shape.

**Resolved: (b).** Stake: JD-1's. Grounded: gRPC encodes "may have
completed" as a property of `DeadlineExceeded`/`Unknown` and picks
`Aborted`/`Unavailable`/`FailedPrecondition` by the client's remedy — one code
per remedy; DDIA's in-doubt 2PC state resolves by querying status before
retry, JD-1's exact stake; REST practice ties control flow to codes, never
message text. Lands in **0026** (Stage 8, as `0026:A8`'s `Applied()` key) and
JDR 0001 §D10's table; **0028** unchanged.

## Interface record

- **JD-1 The `execution_failure` sub-reason carrier.** *(constraint)* — one
  caller-facing slot (`Detail`), one typed applied discriminator
  (`Applied()`, executor-set), `Err` typed for the executor only; no class
  split, no reason enum. Stake: an agent loop deciding whether it must re-read
  before retrying a write, or may retry a read unchanged. *(0026×0028
  Stage-2 joint check; `executor.go::refusalOf`, `model.go::Refusal.Applied`)*

- **JD-2 Invocation-bound argv words.** *(decided)* — `{tag.<key>}` never
  at argv0; a bound value beginning with `-` refuses before spawn; the
  inline-shell promise is unchanged. Stake: the argv a reviewer lints is the
  argv that executes, up to values that cannot be flag-shaped. *(§D2)*
- **JD-3 Exit group of a stale-model write refusal.** *(decided)* — exit 2,
  class unchanged, discriminated by the typed `Err`, one CLI code, reason in
  `findings[]`. Stake: an exit-code-driven agent fixes the model instead of
  retrying it. *(§D3)*
- **JD-4 The applied sense at the envelope.** *(decided)* — a distinct exit-3
  CLI code for the applied write refusal, keyed on `Applied()`; not-applied
  refusals keep `flow-accessor-failed`. Stake: JD-1's, made branchable.
  *(§D4)*

## What this does not decide

- The `carrierDefect`/`Categories()` clause boundary between 0027 and 0028 —
  stays with `cli/0025:C5`, cited by both.
- 0028's reliance on `cli/0016:C4`'s fail-closed read-back reader — 0028's
  instance.
- The concrete `Detail` wording each producer emits — `cli/0026:C1` and
  `0028:C1.3` respectively.
- JDR 0001 §D10's code table and exit codes — cross-referenced, not restated.
  Whether a stale-anchor refusal (`0028:C1.3`: "a stale model, not a missing
  line") belongs in §D10's exit-2 "the model is wrong" group rather than
  `execution_failure`'s exit-3 "repair and re-run unchanged" — filed as §D3
  (open); 0028 carries no §D10 citation today.
- Whether a write command that exits non-zero without a held pipe may have
  mutated — 0004/0025's shipped decision, not reopened here.
