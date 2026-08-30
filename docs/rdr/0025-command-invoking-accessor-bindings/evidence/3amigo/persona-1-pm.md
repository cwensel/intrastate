Model: claude-opus-5[1m]

# Persona 1 — Product Manager

Question: does this RDR actually deliver the user outcome?

Owned starting set: `§problem-statement`, `§approach`, `§decision-rationale`, `0025:MVV`.
**Widened to** `0025:C6`, `0025:A10`, `§capability-dependencies`, `§consequences`,
`§testing-strategy` (S7), `§implementation-plan` phases, and to the source repo
(`internal/cli/flow_state.go`, `internal/cli/flow_resolve.go`,
`internal/cli/flowbind/registry.go`, `internal/table/load.go`).
What sent me there: the Problem Statement's outcome claim is stated in prose
("the caller … must execute the edit itself"), and the MVV's step 3 asserts a
runtime path ("the resolved decision invokes the declared write command") that
the owned sections do not define anywhere. A silence has no line range, so I had
to widen to the plan, the gate, and the code to find whether the outcome lands.

---

## High

### H1 — `0025:MVV` step 3 asserts an outcome the RDR does not build: "the resolved decision invokes the declared write command"

MVV step 3 reads: "drive a state change end to end: **the resolved decision
invokes the declared write command**". `§testing-strategy` S6 restates it as
"the MVV end to end". But no such path exists, and no Phase builds one.

Grounded in source:
- `internal/cli/flow_resolve.go` emits a **plan**, and its own help text tells
  the caller to **transcribe** it by hand: "transcribe these with `--clear
  <key>` on set-state, never `--write key=<clear>`" (`flowResolveExtendedDesc`).
- `internal/cli/flow_state.go:311` is the **sole** non-test call to
  `exec.Write(...)`, inside `flow set-state`, whose writes come from
  `parseWrites(cmd, req.model)` — i.e. from caller-supplied `--write key=value`
  flags, not from a resolved decision. `flow_state.go:273` is explicit that
  `--tag` values "are context only and never become writes (REQ-34)".
- `§implementation-plan` Phases 1–4 build carrier, binding family, registry
  selection, and lint surface. **None** wires resolve→write.

So the delivered outcome is narrower than the MVV's sentence: a caller runs
`flow resolve`, transcribes the plan into `flow set-state --write …` by hand,
and *that* invokes the declared command. Whether that transcription step is in
scope is exactly the decision this ambiguity blocks.

**Decision blocked:** Scope Verification (`0025:G-scope`) cannot be answered —
the gate asks the implementer to confirm the MVV "will be executed during
implementation, not deferred," but as written step 3 either (a) requires an
unbuilt resolve→apply path, making the MVV out of scope for the four Phases, or
(b) means "a caller transcribes the plan into set-state," in which case the
sentence should say so. The implementer cannot pick without a ruling.

### H2 — `§problem-statement`'s stated harm is transcription; the chosen approach removes only the *carrier* half, and the RDR never says the other half survives

The Problem Statement names the harm precisely: "The caller discovers this when
it takes the resolved decision **and must execute the edit itself** — an edit no
model declared and no lint covered, relocating exactly the prose-and-drift
problem the table was built to remove."

Two distinct defects are conjoined there: (i) the edit is not *declared*
(no command carrier), and (ii) the caller is the one *carrying the decision
across* to the edit. This RDR fixes (i). Defect (ii) — the caller hand-copying
`writes{}` from `flow resolve` into `flow set-state --write` — is untouched, and
grounded above it demonstrably persists. `§consequences` claims the positive as
"a decided transition can finally be APPLIED to the artifact the state lives
in," which reads as if (ii) closed too; nothing in `§approach`,
`§consequences`, or `§implementation-plan` records (ii) as surviving or as a
successor's problem. Contrast the Selection LBD in `§load-bearing-decisions`,
which *does* explicitly record a bounded-not-retired residue for the
`verdictFor`/`unreachable` suffixes — the same honesty is absent here.

**Decision blocked:** the Proportionality gate (`0025:G-proportionality`) and
the successor-RDR seed decision. If (ii) is out of scope, it needs naming as a
successor's contract so the product thesis gap is tracked; if in scope, this RDR
owns a second seam (decision→apply wiring) beyond the command carrier, which the
split test would flag.

### H3 — `0025:A10` is `Pending` and `0025:C6` normatively cites a surface that may not ship, leaving the adoption story undecided

`0025:C6` reads `allow_commands = true   # user-scope config (NEW surface,
A10) or --allow-commands`. `0025:A10` Status is **Pending**, Method **Source
Search**, and its Evidence says the positive call is unmade: "whether v1 ships
the flag alone or also the file, and if the file, its search order, its
precedence against the flag, and its absent/malformed behaviour."
`§capability-dependencies` marks the row **Build** with source "**none — no
config subsystem exists**". Grounded: `internal/` has no `config` package
(`accessor`, `cli`, `graphlint`, `guard`, `resolve`, `table`, `version` only).

The PM consequence is not a mechanism gap, it is an **outcome** gap.
`§consequences` sells the friction as "one boolean of adoption friction —
nothing executes until the invoking user sets `allow_commands = true` **once**
in their own config." That "once" is the entire adoption promise. If v1 ships
the flag alone (which `0025:C6` explicitly permits: "The flag alone is a
complete gate for a v1 that ships before the file"), the user outcome is *not*
one-time — every single invocation must carry `--allow-commands`, in every
script and CI step. That is a materially different adoption product, and
`0025:A10`'s own **If wrong** clause concedes it: "the gate degrades to the
flag, which changes the 'one-time opt-in' property the reversal ledger argues
for." `§testing-strategy` S7 also splits on it ("Once A10 fixes the config
surface, the same scenario adds …"), so the test surface is undecided too.

**Decision blocked:** cannot lock. `0025:G-assumptions` requires every critical
assumption Verified; A10 is Pending on a *decision*, not on evidence, so it
cannot be closed by more searching — it is a NEEDS_DECISION fork. Also blocks
the Proportionality/Profile re-derivation, since building a user-scope config
subsystem is a new user-facing surface that likely moves the Profile.

---

## Medium

### M1 — `0025:MVV` end state overstates: "a linted model applies and verifies a real state change through declared commands only"

The end state says "through declared commands **only**". The MVV's own step 1
binds the write to "a declared **wrapper** over `git config` (the A3 spike's
`wrapper-write.sh` shape)", and `§approach` (Integration honesty) concedes the
wrapper "body is script the model does not carry". `§consequences` repeats it as
a Negative. So the proof-of-outcome is: model declares argv → argv names a shell
script the model does not contain → script does the real edit. The reviewable-
authority claim ("the reviewer reads the argv in the model") is true of the
argv, not of what actually mutates the artifact.

This is disclosed honestly three times elsewhere, which is why this is Medium
not High — but the MVV sentence a reader quotes as "did we deliver it" says
"declared commands only", and it is the one sentence that will be cited.

**Decision blocked:** the user-facing claim in docs/release notes, and the
Scope Verification wording — the implementer needs to know whether the MVV
passes with a wrapper doing the write (it must, per step 1) or whether "only"
is a real constraint.

### M2 — `§decision-rationale`'s deciding row "Binds established tools without wrappers" scores O2 = 4, but `0025:A3` establishes writes *always* need a wrapper

The QOC matrix scores O2 at 4 on "Binds established tools without wrappers"
(rationale: "`{artifact}` covers the positional case"), and the prose names that
row as one of the two deciding rows: "it is the only option scoring ≥4 on both
static reviewability and binding established tools — the two clauses the Problem
Statement conjoins."

But `§approach` states, citing the A3 spike: "established-tool **writes**
always take the other form, because a fixed argv cannot carry the planned value
and the tools do not read stdin". The Problem Statement's *primary* gap is the
write side ("Today the second half is missing … the change cannot be APPLIED").
So on the half the problem statement leads with, O2 scores the same as O1 (2 —
"a protocol wrapper per tool"); the 4 is earned entirely on the read/gate half.
O2 still wins the matrix (26 vs 22 with the row at 2 → 24 vs 22), so the
decision holds — but the *reason given* for choosing it misdescribes the outcome
delivered on the problem's leading clause.

**Decision blocked:** the Contradiction Check gate (`0025:G-contradiction`) —
the rationale row and the Approach's integration-honesty paragraph are in
tension about the same fact, and a reviewer re-deriving the choice will hit it.

---

## Low

### L1 — `§approach` cites "the gap both accreted point-fixes (Seam Lineage) patched around", but there is no Seam Lineage section in the record

`§approach` opens with a parenthetical cross-reference "(Seam Lineage)". The
outline has no such section (`§background`, `§technical-environment`,
`§investigation`, `§key-discoveries` are the nearest). `§background` does carry
the substance (kata intrastate#v0hb, RDR 0004 lineage, "No RDR 0001–0024
adjudicates a command-invoking binding"), so the content survives — only the
pointer dangles. A reader trying to confirm the "accreted point-fixes" premise,
which is the motivating narrative for the whole approach, has nowhere to land.

**Decision blocked:** none, strictly — but it weakens the reviewer's ability to
audit the motivating premise during the Finalization Gate sweep.

### L2 — `0025:MVV` step 3's gate wording says "in the invoking configuration" while `0025:C6` says the config may not exist in v1

Step 3: "With `allow_commands = true` **in the invoking configuration** (C6;
without it, the same invocation refuses …)". `0025:C6` and `0025:A10` leave open
that v1 ships `--allow-commands` only. `§testing-strategy` S7 handles this
correctly (flag first, config "once A10 fixes the config surface"); the MVV does
not. Downstream of H3 — resolving A10 resolves this.

**Decision blocked:** none independently; it is a wording follow-on of H3, but
it means the MVV as written is not executable under the flag-only v1.
