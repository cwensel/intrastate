Model: claude-opus-5

# Persona 1 — Product Manager

Owned set read: `§problem-statement`, `§approach`, `§decision-rationale`, `MVV`, plus `§metadata`.

**Widened**, and what sent me: the outcome claim in `§problem-statement` is "zero caller
edits and zero caller transcription", and `MVV` is the only place that outcome is made
executable. To judge whether the MVV actually demonstrates the claimed outcome I had to
read `C1` (esp. C1.6), `§capability-dependencies`, `§existing-infrastructure-audit`,
`§consequences`, `§failure-modes`, `A1`/`A7`/`A9`, and then check the repo: `internal/cli/flow.go`
(`registerTagFlag`), `internal/cli/flow_exec.go` (`buildRequest`/`parseTags`),
`internal/cli/flow_resolve.go`, `internal/accessor/binding.go`.

---

### 1. The MVV's own command cannot bind the tag the C1.6 reader needs on the resolve half

**Anchor:** `0028:MVV` (step 2), with `0028:C1` C1.6 and `0028:A7`.

`MVV` step 2 is the acceptance pipeline and is written as:

```
flow --allow-commands resolve --model … --artifact record=… --artifact readme=… --outcome lock --as json
  | flow --allow-commands set-state --model … --artifact record=… --artifact readme=… --tag nnnn=NNNN --plan -
```

`--tag nnnn=NNNN` appears only on the `set-state` half. But the MVV fixture (paragraph above
step 1) declares "(c) a C1.6 command reader over the README row on role `readme` (`{tag.nnnn}` in
its argv)". C1.6's `binding:` clause says an unbound tag at invocation "refuses `execution_failure`
BEFORE spawn". `flow resolve` runs the narrowed readers (`internal/cli/flow_resolve.go` header:
"parse input -> load model -> run the narrowed readers -> ONE kernel call"), and it parses tags —
`internal/cli/flow_exec.go::buildRequest(cmd, withTags)` calls `parseTags`, and
`internal/cli/flow_resolve.go:120` registers the flag. So if the readme reader is narrowed in
(its key is referenced by the row being resolved, which it must be for step 3's
"both confirmed by read-back"), the pipeline as written refuses at the *resolve* stage, before
`set-state` is ever reached — the MVV fails on its own happy path for a reason that is not the
carrier under test.

`A7` verified `--tag` only on `set-state` ("so `{tag.nnnn}` is bindable on the invocation"); it
never checks the resolve half, and the `If wrong` branch discusses envelope linkage rather than
this omission. Either the command needs `--tag nnnn=NNNN` on both halves, or the record must
state that the readme reader is *not* narrowed into resolve and the readme status is read only
at read-back — in which case step 3's "both confirmed by read-back" is doing work the record
never explains.

**Blocks:** whether the acceptance scenario is executable as written — an implementer running
`MVV` step 2 verbatim gets a refusal and cannot tell whether the carrier is wrong or the
scenario is. It also blocks the Stage-8 exit criterion, since `MVV` is the record's only
end-to-end proof of the stated user outcome.

---

### 2. "Zero caller edits" is asserted as the outcome but the record's own dependency table says the consumer must add a verb

**Anchor:** `0028:§problem-statement`, `0028:§capability-dependencies` (last row), `0028:MVV`.

`§problem-statement` states the outcome as "a linted table to decide a state change and
`intrastate` to apply it to the record the state lives in, with **zero caller edits and zero
caller transcription**", and `§decision-rationale` restates it as "nothing the caller types and
nothing lint cannot see". But `§capability-dependencies`' final row — "Row-addressed command
reader over a shared artifact (`{tag.<key>}` in argv) … **the consumer adds the verb**" — and
`MVV` step 2's `--tag nnnn=NNNN` are both caller-side work: a new projector invocation form the
consumer must build, plus a per-invocation tag the caller must type with the correct record
number.

The record never reconciles these. "Zero caller transcription" plausibly means "the caller never
types the planned VALUE" (which A7 does establish), and the tag is an identity rather than a
value — but the RDR never says so. As written, the headline outcome and the dependency table
contradict each other, and a reader cannot tell which of the two is the promise being shipped.

**Blocks:** the adoption/acceptance decision — whether "done" means the consumer's
`rdr-write.toml` works with no consumer-side change (false, per the dependency table) or means
"no wrapper script and no typed value" (true). Also blocks sizing the consumer-side work, which
is currently invisible in `§implementation-plan`'s four phases.

---

### 3. The README half of the user outcome is demonstrated only against a fixture that excludes the live data

**Anchor:** `0028:MVV` (fixture paragraph), `0028:F6`, `0028:§key-discoveries`.

The MVV fixture explicitly scopes to "a README with all three index rows in the linked
`| [NNNN](NNNN-slug.md) |` form", and then states: "the unlinked `| NNNN |` form two live rows
carry is out of the fixture and named in Failure Modes". `F6` and `§key-discoveries` confirm the
real consumer README carries rows 0001 and 0006 in that unlinked form and that they refuse
`edit_anchor_unmatched`.

The record disposes of this as "pre-existing data drift, not a defect of this carrier", and that
reasoning holds for the *contract*. But from the outcome side the consequence is unstated: the
motivating user (rdr#yjye) cannot flip the state of records 0001 and 0006 through this feature
until someone normalizes the README, and no owner, ticket, or ordering for that normalization
appears anywhere in the record — not in `§implementation-plan`, `§prerequisites`, or `§metadata`
Related Issues. The MVV therefore passes while the shipped feature does not deliver the outcome
for two of the consumer's live rows.

**Blocks:** the ship/no-ship call and the definition of "delivered" for the consumer. An
implementer can pass every MVV step and still hand back a feature the consumer cannot use on
its real README. A named prerequisite (normalize the two rows) or an explicit
"partial adoption, N rows deferred" consequence would resolve it.

---

### 4. `§approach` and `§decision-rationale` sell C1.6 as a small companion; the record's own evidence makes it a co-equal half of the outcome

**Anchor:** `0028:§approach` (final paragraph), `0028:§decision-rationale`, `0028:C1` (C1.6),
`0028:A1`.

`§approach` closes with "**One companion clause (C1.6)** lets the same declared context tag
row-address a command READER", framing it as an addendum to the `edit` carrier. `§decision-rationale`'s
scored matrix (drift removal / blast radius / prior-art alignment) is entirely about the write
carrier and never scores C1.6 or its alternatives at all — the only justification offered is
one premortem line ("the README half of the scenario was unverifiable as briefed (P-15) — folded
as C1.6").

But the record's own material makes C1.6 load-bearing for the *stated outcome*, not incidental:
`§approach` itself admits "without it the README half of the consumer scenario is unverifiable
by construction"; `A1` shows the seam extension exists *because* a command reader needs the same
tags (`Read(ctx, art Artifact, requested []string)` — verified on disk in
`internal/accessor/binding.go`, the same `art` value); it extends 0025:C2's *closed* vocabulary
(a declared `Overrides` in `§metadata`); and it turns a refusal into an acceptance
(`§decision-rationale` Arm 3). No alternative to C1.6 is ever weighed — e.g. a per-record role
binding, which `A1`'s `If wrong` names as the fallback but which is never scored.

**Blocks:** the decision of whether C1.6 belongs in this RDR at all, and — via `G-proportionality`,
which asks the profile to be re-derived from contract count — whether this record is the sole
author of one load-bearing contract or is quietly locking two seams (a write carrier and an
argv-vocabulary extension) under one `C1` id.

---

### 5. `§decision-rationale`'s process ledger crowds out the rationale a reader needs

**Anchor:** `0028:§decision-rationale` (the `Premortem: / Ground-sweep: / Joint-check:` block).

Roughly half of `§decision-rationale`'s bytes are a single unbroken paragraph of pre-lock
bookkeeping — joint-check arms, grepped placeholder literals, per-peer dispositions, an inline
citation correction. The actual rationale (the three matrix rows, the one-line rejections) is
crisp and precedes it. The ledger is process residue, not the reasoning behind the choice, and
it is in the section a future reader opens to ask "why this shape?".

This is the mildest finding and is a proportionality/readability issue rather than a defect in
the decision. But `G-proportionality` explicitly asks whether the document is right-sized and
whether provenance prose should be stripped, and this block is the clearest candidate — the
same content already lives in `evidence/joint-check/arms.md`, which the paragraph cites.

**Blocks:** nothing in the design; it blocks the `G-proportionality` gate answer, which cannot
honestly say "right-sized" without disposing of this block.

---

## Not findings (checked and held)

- `A1`'s seam-extension claim holds on disk: `internal/accessor/binding.go` declares
  `WriteBinding.Apply(ctx, art Artifact, planned []resolve.Tag) error` and
  `ReadBinding.Read(ctx, art Artifact, requested []string)` — the same `art`, and
  `Artifact` carries no context field. One context map on `Artifact` serves both, as claimed.
- `A7`'s `registerTagFlag` claim holds: `internal/cli/flow.go:182`, wired at
  `internal/cli/flow_state.go:199`.
- The `§load-bearing-decisions` `D-selection-predicate` claim "no non-test `regexp` import under
  `internal/`" holds — `rg -l '"regexp"' internal/ --type go` excluding tests returns nothing.
- `§alternatives-considered` weighs A–E on stated criteria with per-option rejection reasons; the
  rejections of B and E trace directly to the problem statement's drift framing. No defect.
