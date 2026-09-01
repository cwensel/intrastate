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
  `0004`'s closed set, and `0028:C3`'s own "no new refusal class"; it also
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
gain the `Applied()` derivation (rule 3). Lands in **0028** — `0028:C3`'s
`order:` line drops `applied: false` from `Detail` (rule id and reason token
remain); `0028:A8` is confirmed as a constraint of this entry, and its
If-wrong "owned here" is void — a typed field on the carrier is this
registry's, not 0028's. **0004** and **0025** are Implemented and not edited;
0026 and 0028 are the successors at their seams and cite this entry.
**0027** is bound as a member; nothing lands in it.

## Interface record

- **JD-1 The `execution_failure` sub-reason carrier.** *(constraint)* — one
  caller-facing slot (`Detail`), one typed applied discriminator
  (`Applied()`, executor-set), `Err` typed for the executor only; no class
  split, no reason enum. Stake: an agent loop deciding whether it must re-read
  before retrying a write, or may retry a read unchanged. *(0026×0028
  Stage-2 joint check; `executor.go::refusalOf`, `model.go::Refusal.Applied`)*

## What this does not decide

- The `carrierDefect`/`Categories()` clause boundary between 0027 and 0028 —
  stays with `cli/0025:C5`, cited by both.
- 0028's reliance on `cli/0016:C4`'s fail-closed read-back reader — 0028's
  instance.
- The concrete `Detail` wording each producer emits — `cli/0026:C1` and
  `0028:C3` respectively.
- JDR 0001 §D10's code table and exit codes — cross-referenced, not restated.
  Whether a stale-anchor refusal (`0028:C3`: "a stale model, not a missing
  line") belongs in §D10's exit-2 "the model is wrong" group rather than
  `execution_failure`'s exit-3 "repair and re-run unchanged" is §D10's
  question; 0028 answers it there by citation.
- Whether a write command that exits non-zero without a held pipe may have
  mutated — 0004/0025's shipped decision, not reopened here.
